# Transfer data from GCP Redis to Memorystore Valkey

This example demonstrates how to migrate data from a Redis (Memorystore for Redis) instance to a Valkey (Memorystore for Valkey) instance using Cloud Storage as an intermediate.

You will:

1. Create a VPC network and a Cloud Storage bucket.
2. Prepare a local `dump.rdb` file and upload it to the bucket.
3. Create a Redis instance and import the uploaded data.
4. Export a dump from the Redis instance to the bucket.
5. Create a new Valkey instance by importing the dump from the bucket.
6. Create a backup of the Valkey instance and export it to the bucket.

The directory contains helper scripts to perform each step.

## Prerequisites

- A GCP project (`gcloud config set project <PROJECT_ID>` is already done)
- IAM roles such as:
  - Compute Network Admin (for VPC and subnets)
  - Storage Admin (for Cloud Storage)
  - Memorystore Admin (for Redis/Valkey)
  - Service Account Admin (if using custom service accounts)
- `gcloud` CLI installed
- Docker if you want to inspect Redis data locally

## 1. Create VPC network and Cloud Storage bucket

Use the `create-vpc-and-bucket.sh` script to create a custom VPC network, subnet, and a GCS bucket.

```sh
./create-vpc-and-bucket.sh
```

The script is expected to:

- Create a VPC (for example: `${NETWORK_NAME}`)
- Create a subnet in `${REGION}`
- Create a GCS bucket (for example: `gs://${BUCKET_NAME}`)

Check and update the script variables if you want to use different names.

## 2. Prepare a local `dump.rdb` and upload it to the bucket

Prepare a Redis RDB dump file locally. You can either:

- Use an existing Redis environment and create a dump, or
- Use the sample file in `data/sample-data-dump.rdb` as a starting point.

Upload the RDB file to the bucket.

If you have your own `dump.rdb`, upload that instead.

## 3. Create a Redis instance and import the uploaded data

Use the `create-redis.sh` script to create a Memorystore for Redis instance and import the RDB file from Cloud Storage.

```sh
./create-redis.sh
```

## 4. Export dump data from the Redis instance to the bucket

After the Redis instance is up and populated, export its data again to Cloud Storage.  

## 5. Create a new Valkey instance importing the dump data from the bucket

Use the `create-valkey-from-gcs.sh` script to create a Memorystore for Valkey instance, importing the dump file from the bucket.

```sh
./create-valkey-from-gcs.sh
```

The script is expected to:

- Create a Valkey instance in `${REGION}`
- Attach it to the same VPC/subnet
- Import the RDB file from `gs://${BUCKET_NAME}` (e.g., from `redis-export/dump.rdb`)

Again, verify the bucket name, region, and object path inside the script.

## 6. Create a backup of the Valkey instance and export it to the bucket

Use the `backup-valkey-into-gcs.sh` script to create a backup of the Valkey instance and export it to Cloud Storage.

```sh
./backup-valkey-into-gcs.sh
```

This should:

- Create a backup of the Valkey instance
- Export the backup RDB file to `gs://${BUCKET_NAME}` (for example, `backup-valkey-memorystore-transfer-data.rdb`)

The backup file can be found in the `data/` directory locally or in the bucket depending on how the script is implemented.

## Data samples

The `data/` directory contains the following sample files:

- `sample-data-dump.rdb`: Sample dump with example key/value data
- `backup-valkey-memorystore-transfer-data.rdb`: Example backup of a Valkey instance

You can use these files to quickly test the transfer flow without preparing your own data.

## Verification and troubleshooting

### Verify data using local Docker and `docker-compose.yaml`

This directory contains a `docker-compose.yaml` that starts a local Redis container and loads an RDB file so that you can inspect the data.

1. **Start local Redis with the RDB file you want to inspect**

   By default, `docker-compose.yaml` mounts the `data/` directory and loads an RDB file named `dump.rdb`.

   ```sh
   docker compose up -d
   ```

2. **Connect with `redis-cli` and inspect keys/values**

   ```sh
   docker compose exec redis redis-cli
   ```

   In the `redis-cli` shell, run commands such as:

   ```text
   KEYS *
   GET some:key
   TYPE some:zset
   ZRANGE some:zset 0 -1 WITHSCORES
   ```

   Compare the data before and after each migration step (e.g., original dump, Redis export, Valkey backup) by switching which RDB file is mounted in `docker-compose.yaml` or by copying different dump files into the `data/` directory.
