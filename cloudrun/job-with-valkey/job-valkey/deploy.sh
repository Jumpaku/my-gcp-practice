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

gcloud compute networks describe vpc-job-with-valkey --project="${PROJECT_ID}" || \
  gcloud compute networks create vpc-job-with-valkey \
    --project="${PROJECT_ID}" \
    --subnet-mode=custom \
    --mtu=1460 \
    --bgp-routing-mode=regional \
    --bgp-best-path-selection-mode=legacy
gcloud compute networks subnets describe subnet-job-with-valkey --project="${PROJECT_ID}" --region="${REGION}" || \
  gcloud compute networks subnets create subnet-job-with-valkey \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --range=10.0.0.0/24 \
    --stack-type=IPV4_ONLY \
    --network=vpc-job-with-valkey
gcloud network-connectivity service-connection-policies describe psc-job-with-valkey --project="${PROJECT_ID}" --region="${REGION}" || \
  gcloud network-connectivity service-connection-policies create psc-job-with-valkey \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --network=vpc-job-with-valkey \
    --service-class=gcp-memorystore \
    --psc-connection-limit=5 \
    --subnets="https://www.googleapis.com/compute/v1/projects/${PROJECT_ID}/regions/${REGION}/subnetworks/subnet-job-with-valkey"

gcloud memorystore instances describe valkey-job-with-valkey --project="${PROJECT_ID}" --location="${REGION}" || \
  gcloud memorystore instances create valkey-job-with-valkey \
    --project="${PROJECT_ID}" \
    --location="${REGION}" \
    --endpoints=connections="[{pscAutoConnection={network=projects/${PROJECT_ID}/global/networks/vpc-job-with-valkey,port=6379,projectId=${PROJECT_ID}}}]" \
    --shard-count=1 \
    --node-type=standard-small \
    --replica-count=0 \
    --engine-version=VALKEY_8_0

gcloud run jobs deploy "${JOB}" \
  --project "${PROJECT_ID}" \
  --region "${REGION}" \
  --service-account="${SERVICE_ACCOUNT}" \
  --max-retries=0 \
  --tasks=1 \
  --source . \
  --network vpc-job-with-valkey \
  --subnet subnet-job-with-valkey \
  --vpc-egress private-ranges-only
