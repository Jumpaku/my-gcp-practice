#!/bin/sh

set -eux

gcloud config list
gcloud auth list

echo "$*"
ls
pwd
uname -a
env
date

MAX_INSTANCES="$(TZ="Asia/Tokyo" date +'%M')"
gcloud run deploy "${TARGET_SERVICE}" \
  --region=asia-northeast1 \
  --image=us-docker.pkg.dev/cloudrun/container/hello \
  --max="${MAX_INSTANCES}"