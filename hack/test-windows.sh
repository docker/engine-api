#!bash -x
mkdir -p "/tmp/$BUILD_TAG"
trap "rm -rf /tmp/$BUILD_TAG" EXIT

export GOVERSION=1.5.3
export GOPATH="/tmp/$BUILD_TAG"
mkdir -p "$GOPATH/src/$GOPACKAGE"
rm -rf "$GOPATH/src/$GOPACKAGE"
cp -R "$PWD" "$GOPATH/src/$GOPACKAGE"

export GOBIN="$GOPATH/bin"
export PATH="$PATH:$GOBIN"
export GOCOV_ARGS="-race -tags=test"

mkdir -p "$GOBIN"
curl -sSLo "$GOBIN/golang_test.sh" "https://${GITHUB_TOKEN}@raw.githubusercontent.com/docker/tools-team/master/Dockerfiles/golang-tester/entrypoint.sh"
chmod u+x "$GOBIN/golang_test.sh"

## Not using gimme currently
# gimme "$GOVERSION"
# curl -sSLo "$GOBIN/gimme "https://${GITHUB_TOKEN}@raw.githubusercontent.com/travis-ci/gimme/master/gimme
# chmod u+x "$GOBIN/gimme"
cp "$(which true)" "$GOBIN/gimme"

go version

go get github.com/golang/lint/golint
go get github.com/fzipp/gocyclo
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml
go get bitbucket.org/tebeka/go2xunit

mkdir -p results
export OUTPUT_DIR="results"
exec bash -x golang_test.sh
