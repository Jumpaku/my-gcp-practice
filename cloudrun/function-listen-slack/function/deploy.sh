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

if [ "${SERVICE}" = "" ]; then
  echo "Please set SERVICE"
  exit 1
fi

if [ "${SERVICE_ACCOUNT}" = "" ]; then
  echo "Please set SERVICE_ACCOUNT"
  exit 1
fi

if [ "${SLACKAPP_SIGNING_SECRET}" = "" ]; then
  echo "Please set SLACKAPP_SIGNING_SECRET"
  exit 1
fi

gcloud functions deploy "${SERVICE}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --gen2 \
  --trigger-http \
  --allow-unauthenticated \
  --service-account="${SERVICE_ACCOUNT}" \
  --build-service-account="projects/${PROJECT_ID}/serviceAccounts/${SERVICE_ACCOUNT}" \
  --set-secrets="SLACK_SIGNING_SECRET=projects/${PROJECT_ID}/secrets/${SLACKAPP_SIGNING_SECRET}/versions/latest" \
  --runtime go125 \
  --source . \
  --entry-point HandleSlackEvent
