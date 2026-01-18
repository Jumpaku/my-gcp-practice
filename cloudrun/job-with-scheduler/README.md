# Cloud Run Job with Cloud Scheduler

This sample describes an architecture where Cloud Scheduler periodically triggers a Cloud Run Job, and the result is verified via a sample Cloud Run service.

## Prerequisites

- You have a GCP project (`gcloud config set project <PROJECT_ID>` is already done)
- You have IAM roles such as Cloud Run Admin / Cloud Scheduler Admin / Service Account Admin
- `gcloud` CLI is installed

## Overall architecture

1. Sample Cloud Run service (HTTP)
2. Cloud Run Job
3. Cloud Scheduler (Pub/Sub or HTTP)

## GCP resources to prepare

1. **Deploy the sample Cloud Run service**
   You can either deploy a simple "Hello" service, or reuse an existing Cloud Run service.  
   In this example, we use the `gcr.io/cloudrun/hello` image.

2. **Deploy the Cloud Run Job**
   
3. **Configure Cloud Scheduler**


Use the script in the `job-with-scheduler/job-schedule/` directory to prepare the above.

```sh
cd cloudrun/job-with-scheduler/job-schedule
./deploy.sh
```

## How to verify the behavior

1. In the Cloud Scheduler console, click **Run now** and confirm that the Cloud Run Job starts successfully.
2. Check the Cloud Run Job execution history and logs, and confirm that the calls to the sample service succeed.
3. Check the logs and responses on the sample service side and confirm that it is being invoked as expected.

## References

- [Cloud Run](https://cloud.google.com/run/docs)
- [Cloud Run Jobs](https://cloud.google.com/run/docs/jobs)
- [Cloud Scheduler](https://cloud.google.com/scheduler/docs)
