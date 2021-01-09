#!/bin/bash

set -euxo pipefail

cd $(dirname $0)

CGO_ENABLE=0 GOOS=linux GOARCH=arm64 go build -o pro_deploy.arm64 .

echo "编译成功 🎉🎉🎉  "
