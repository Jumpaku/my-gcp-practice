# Cloud Run Job with Valkey (MemoryStore)

This sample shows how to connect from a Cloud Run Job to Valkey (Redis-compatible MemoryStore).

## Prerequisites

- You have a GCP project (`gcloud config set project <PROJECT_ID>` is already done)
- You have IAM roles such as VPC Admin / Compute Network Admin / Memorystore Admin / Cloud Run Admin / Service Account Admin
- `gcloud` CLI is installed

## Overall architecture

1. VPC & subnet
2. Create an endpoint to Valkey via Private Service Connect (PSC)
3. Valkey instance
4. Cloud Run Job (connects to Valkey via a VPC connector)

## GCP resources to prepare

1. Create a VPC and subnet

2. Create a Serverless VPC Access connector

3. Create a Valkey (MemoryStore) instance
    You can create the instance either from the GCP Console or using the CLI. Here is a CLI example:
    After creation, note down the connection information (host name and port).

## Deploy the Cloud Run Job

Deploy the Go application in the `job-with-valkey/job-valkey/` directory so that it connects to Valkey via the VPC connector.

```sh
cd job-valkey
./deploy.sh
```

The `deploy.sh` script is expected to configure settings such as the following:

- Specify the Serverless VPC Access connector with `--vpc-connector=${CONNECTOR_NAME}`
- Set the Valkey host and port as environment variables (for example: `VALKEY_HOST`, `VALKEY_PORT`)
- Set the execution service account to `cr-job-valkey-sa`

Update `deploy.sh` as needed to reflect your Valkey connection information and VPC connector name.

## How to verify the behavior

1. Execute the job
2. Check Cloud Logging to confirm that operations such as SET/GET against Valkey are succeeding.
