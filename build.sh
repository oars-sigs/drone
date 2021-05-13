#!/bin/bash
set -e

go mod vendor
cp -f ui/dist/dist_gen.go vendor/github.com/drone/drone-ui/dist/
go build -ldflags "-extldflags \"-static\"" -mod vendor -o bin/drone-server github.com/oars-sigs/drone/cmd/drone-server