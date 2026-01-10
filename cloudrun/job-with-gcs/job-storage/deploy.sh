#!/bin/sh

set -eux

if [ "${PROJECT_ID}" = "" ]; then
  echo "Please set PROJECT_ID"
  exit 1
fi

if [ "${REGION}" = "" ]; then
  echo "Please set REGION"
  exit 1
fi

if [ "${JOB}" = "" ]; then
  echo "Please set JOB"
  exit 1
fi

if [ "${SERVICE_ACCOUNT}" = "" ]; then
  echo "Please set SERVICE_ACCOUNT"
  exit 1
fi

if [ "${GCS_BUCKET}" = "" ]; then
  echo "Please set GCS_BUCKET"
  exit 1
fi

gcloud storage buckets describe "gs://${GCS_BUCKET}" \
  || gcloud storage buckets create "gs://${GCS_BUCKET}" \
    --project="${PROJECT_ID}" \
    --location="${REGION}" \
    --uniform-bucket-level-access \
    --enable-hierarchical-namespace


gcloud run jobs deploy "${JOB}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --service-account="${SERVICE_ACCOUNT}" \
  --max-retries 0 \
  --tasks=1 \
  --source . \
  --add-volume name=volume_job-with-gcs,type=cloud-storage,bucket="${GCS_BUCKET}" \
  --add-volume-mount volume=volume_job-with-gcs,mount-path=/workspace/gcs_storage
