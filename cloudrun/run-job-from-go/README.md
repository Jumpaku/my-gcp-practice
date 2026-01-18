# Run Cloud Run Job from Go

This sample shows how to start a Cloud Run Job from a Go client.

## Prerequisites

- You have a GCP project (`gcloud config set project <PROJECT_ID>` is already done)
- You have IAM roles such as Cloud Run Admin / Service Account Admin / Cloud Run Invoker
- `gcloud` CLI and Go 1.22+ are installed

## GCP resources to prepare

Deploy the sample job from the `job/` directory.

```sh
cd ./job
./deploy.sh
```

Inside `deploy.sh`, configure the job name (for example: `sample-job`), region, container image, and execution service account.  
After deployment, note down the job name and region.

## Running the Go client

From `run-job-from-go/main.go`, the program calls the Cloud Run Jobs execution API.

1. Set the environment variables used in `main.go`:

   - `CLOUD_RUN_JOB_NAME`: Name of the job to execute
   - `CLOUD_RUN_JOB_REGION`: Region where the job is deployed
   - `GOOGLE_CLOUD_PROJECT`: Project ID

   ```sh
   export CLOUD_RUN_JOB_NAME=sample-job
   export CLOUD_RUN_JOB_REGION=asia-northeast1
   export GOOGLE_CLOUD_PROJECT=$(gcloud config get-value project)
   ```

2. Run the Go client:

   ```sh
   cd cloudrun/run-job-from-go
   go run main.go
   ```

   If the call succeeds, a new execution of the specified Cloud Run Job will be created, and the execution ID and status will be printed in the logs.

## Verifying the behavior

- In the Cloud Console, go to **Cloud Run > Jobs** and confirm that a new entry has been added to the job's execution history.
- Check the logs for each execution and verify that the job is performing the expected processing.
