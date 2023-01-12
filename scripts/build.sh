#!/usr/bin/env bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. ${SCRIPT_DIR}/common.sh

CGO_ENABLED=0 go build -ldflags " -s -w -X 'main.Version="${VERSION}"' -X 'main.GitRev="${GIT_REV}"'"