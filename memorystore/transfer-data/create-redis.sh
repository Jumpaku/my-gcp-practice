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

gcloud redis instances create redis-memorystore-transfer-data \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --network="projects/${PROJECT_ID}/global/networks/vpc-memorystore-transfer-data" \
  --connect-mode=DIRECT_PEERING \
  --tier=basic \
  --size=1 \
  --redis-version=redis_7_2