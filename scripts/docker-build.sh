#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. ${SCRIPT_DIR}/common.sh

REPO=lalyos/bookline

docker build -t ${REPO} .
docker tag ${REPO} ${REPO}:${VERSION}