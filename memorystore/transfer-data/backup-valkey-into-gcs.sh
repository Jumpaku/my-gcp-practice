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

BACKUP_ID="backup-valkey-memorystore-transfer-data_$(date +%Y%m%d-%H%M%S)"
gcloud memorystore instances backup valkey-memorystore-transfer-data \
  --project="${PROJECT_ID}" \
  --location="${REGION}" \
  --backup-id="${BACKUP_ID}"

BACKUP_COLLECTION=$(
  gcloud memorystore backup-collections list \
    --location="${REGION}" \
    --project="${PROJECT_ID}" \
    --filter="instance=projects/${PROJECT_ID}/locations/${REGION}/instances/valkey-memorystore-transfer-data" \
    --format='value(name)'
)

gcloud memorystore backup-collections backups export "${BACKUP_COLLECTION}/backups/${BACKUP_ID}" \
  --location="${REGION}" \
  --project="${PROJECT_ID}" \
  --gcs-bucket="${GCS_BUCKET}"
