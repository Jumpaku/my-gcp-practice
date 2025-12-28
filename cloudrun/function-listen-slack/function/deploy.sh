#!/bin/sh

set -eux

PROJECT_ID="my-gcp-practice"
REGION="asia-northeast1"
SERVICE="function-listen-slack"
SERVICE_ACCOUNT='cloud-run-worker@my-gcp-practice.iam.gserviceaccount.com'
SLACKAPP_SIGNING_SECRET=SLACKAPP_DEMO_SIGNING_SECRET

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
