#!/bin/bash -x
docker pull dockerautomation/golang-tester:gimme
docker run \
--rm \
-v "$(pwd):/go/src/$GOPACKAGE" \
-v "$(pwd)/results:/output" \
-e "GOVERSION=$GOVERSION" \
-e "GOCOV_ARGS=-race -tags=test" \
-e "GOCYCLO_MAX" \
-e "GOPACKAGE" \
dockerautomation/golang-tester:gimme
