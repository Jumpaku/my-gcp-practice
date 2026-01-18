#!/bin/sh

set -eux

if [ -z "${PROJECT_ID}" ]; then
  echo "Please set PROJECT_ID"
  exit 1
fi

if [ -z "${REGION}" ]; then
  echo "Please set REGION"
  exit 1
fi

if [ -z "${GCS_BUCKET}" ]; then
  echo "Please set GCS_BUCKET"
  exit 1
fi

gcloud memorystore instances create valkey-memorystore-transfer-data \
  --project="${PROJECT_ID}" \
  --location="${REGION}" \
  --endpoints=connections="[{pscAutoConnection={network=projects/${PROJECT_ID}/global/networks/vpc-memorystore-transfer-data,port=6379,projectId=${PROJECT_ID}}}]" \
  --shard-count=1 \
  --node-type=standard-small \
  --replica-count=0 \
  --engine-version=VALKEY_8_0 \
  --gcs-source-uris="gs://${GCS_BUCKET}/sample-data-dump.rdb"
