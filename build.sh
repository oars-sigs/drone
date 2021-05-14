#!/bin/bash
set -e

go mod vendor
cp -f ui/dist/dist_gen.go vendor/github.com/drone/drone-ui/dist/
sed -i "89i 	case scm.Driver(8):" vendor/github.com/drone/drone/service/netrc/netrc.go
sed -i "90i 			netrc.Login = \"oauth2\"" vendor/github.com/drone/drone/service/netrc/netrc.go
sed -i "91i 			netrc.Password = user.Token" vendor/github.com/drone/drone/service/netrc/netrc.go
go build -ldflags "-extldflags \"-static\"" -mod vendor -o bin/drone-server github.com/oars-sigs/drone/cmd/drone-server