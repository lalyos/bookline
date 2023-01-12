#!/usr/bin/env bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
. ${SCRIPT_DIR}/common.sh

go build -ldflags "-X 'main.Version="${VERSION}"' -X 'main.GitRev="${GIT_REV}"'"