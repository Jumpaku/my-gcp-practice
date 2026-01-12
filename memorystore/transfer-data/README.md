
Example of transferring data from GCP Redis into GCP Memorystore Valkey.

1. Create VPC and Bucket.
2. Prepare a redis dump.rdb file at local and upload to the bucket.
3. Create a Redis instance and import the uploaded data.
4. Export dump data from the Redis instance into bucket.
5. Create a new Valkey instance importing the dump data from bucket.
6. Create a backup of the Valkey instance and export it into bucket.
