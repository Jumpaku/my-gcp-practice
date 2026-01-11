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

if [ "${SERVICE}" = "" ]; then
  echo "Please set SERVICE"
  exit 1
fi

if [ "${SERVICE_ACCOUNT}" = "" ]; then
  echo "Please set SERVICE_ACCOUNT"
  exit 1
fi

gcloud run deploy "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --image=us-docker.pkg.dev/cloudrun/container/hello \
  --memory=128Mi \
  --cpu=1 \
  --no-allow-unauthenticated

gcloud run jobs deploy "${JOB}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --service-account="${SERVICE_ACCOUNT}" \
  --max-retries 0 \
  --tasks=1 \
  --source . \
  --set-env-vars="TARGET_SERVICE=${SERVICE}"

gcloud scheduler jobs describe --project "${PROJECT_ID}" --location="${REGION}" "scheduler-job-with-scheduler" \
  || gcloud scheduler jobs create http "scheduler-job-with-scheduler" \
    --project "${PROJECT_ID}" \
    --location="${REGION}" \
    --schedule="0 0 1 1 *" \
    --time-zone='Asia/Tokyo' \
    --uri="https://${REGION}-run.googleapis.com/v2/projects/${PROJECT_ID}/locations/${REGION}/jobs/${JOB}:run" \
    --http-method=POST \
    --oauth-service-account-email="${SERVICE_ACCOUNT}"
