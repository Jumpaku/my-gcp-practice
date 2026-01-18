# Cloud Run Job with GCS

This sample demonstrates how to access and manipulate a Cloud Storage bucket from a Cloud Run Job.

## Prerequisites

- You have a GCP project (`gcloud config set project <PROJECT_ID>` is already done)
- You have IAM roles such as Cloud Run Admin / Storage Admin / Service Account Admin
- `gcloud` CLI is installed

## GCP resources to prepare

1. **Create a Cloud Storage bucket**
   Create the bucket that the job will read from and write to.

2. **Create a service account for the Cloud Run Job**

3. **Deploy the Cloud Run Job**

Use the script in the `job-with-gcs/job-storage/` directory to create the job.

```sh
cd cloudrun/job-with-gcs/job-storage
./deploy.sh
```

The `deploy.sh` script is expected to configure at least the following:

- Pass `BUCKET_NAME` to the job via `--set-env-vars`
- Set the execution service account to `cr-job-gcs-sa`

## Run the job and verify behavior

1. Manually execute the job

2. Check the Cloud Storage bucket
   Confirm that the files written by the job exist in the bucket.

3. Check the job logs in Cloud Logging and verify that the GCS API calls have succeeded.
