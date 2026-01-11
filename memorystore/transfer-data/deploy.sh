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

gcloud compute networks describe vpc-memorystore-transfer-data --project="${PROJECT_ID}" || \
  gcloud compute networks create vpc-memorystore-transfer-data \
    --project="${PROJECT_ID}" \
    --subnet-mode=custom \
    --mtu=1460 \
    --bgp-routing-mode=regional \
    --bgp-best-path-selection-mode=legacy
gcloud compute networks subnets describe subnet-memorystore-transfer-data --project="${PROJECT_ID}" --region="${REGION}" || \
  gcloud compute networks subnets create subnet-memorystore-transfer-data \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --range=10.0.0.0/27 \
    --stack-type=IPV4_ONLY \
    --network=vpc-memorystore-transfer-data
gcloud network-connectivity service-connection-policies describe psc-memorystore-transfer-data --project="${PROJECT_ID}" --region="${REGION}" || \
  gcloud network-connectivity service-connection-policies create psc-memorystore-transfer-data \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --network=vpc-memorystore-transfer-data \
    --service-class=gcp-memorystore \
    --psc-connection-limit=5 \
    --subnets="https://www.googleapis.com/compute/v1/projects/${PROJECT_ID}/regions/${REGION}/subnetworks/subnet-memorystore-transfer-data"


gcloud storage buckets describe "gs://${GCS_BUCKET}" \
  || gcloud storage buckets create "gs://${GCS_BUCKET}" \
    --project="${PROJECT_ID}" \
    --location="${REGION}" \
    --uniform-bucket-level-access


gcloud redis instances describe redis-memorystore-transfer-data --project="${PROJECT_ID}" --region="${REGION}" || \
  gcloud redis instances create redis-memorystore-transfer-data \
    --project="${PROJECT_ID}" \
    --region="${REGION}" \
    --network="projects/${PROJECT_ID}/global/networks/vpc-memorystore-transfer-data" \
    --connect-mode=DIRECT_PEERING \
    --tier=basic \
    --size=1 \
    --redis-version=redis_7_2


gcloud memorystore instances describe valkey-memorystore-transfer-data --project="${PROJECT_ID}" --location="${REGION}" || \
  gcloud memorystore instances create valkey-memorystore-transfer-data \
    --project="${PROJECT_ID}" \
    --location="${REGION}" \
    --endpoints=connections="[{pscAutoConnection={network=projects/${PROJECT_ID}/global/networks/vpc-memorystore-transfer-data,port=6379,projectId=${PROJECT_ID}}}]" \
    --shard-count=1 \
    --node-type=standard-small \
    --replica-count=0 \
    --engine-version=VALKEY_8_0

