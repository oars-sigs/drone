#!/bin/bash
set -e

go build -o bin/drone-server github.com/oars-sigs/drone/cmd/drone-server 
docker build -t registry.cn-shenzhen.aliyuncs.com/oars/drone -f Dockerfile .
docker push registry.cn-shenzhen.aliyuncs.com/oars/drone