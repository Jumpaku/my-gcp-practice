#!/bin/sh

set -eux

(
  echo "$*"
  env
  ls
  pwd
  uname -a
) > /workspace/gcs_storage/stdout.txt 2> /workspace/gcs_storage/stderr.txt

