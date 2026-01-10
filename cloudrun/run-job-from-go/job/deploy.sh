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

gcloud run jobs deploy "${JOB}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --service-account="${SERVICE_ACCOUNT}" \
  --max-retries 0 \
  --tasks=1 \
  --source .
