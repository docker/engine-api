.PHONY: all deps test validate

all: deps test validate

deps:
	go get -t ./...
	go get github.com/golang/lint/golint

test:
	go test -tags=test ./...

test_experimental:
	go test -tags="test experimental" ./...

validate:
	go vet ./...
	test -z "$(golint ./... | tee /dev/stderr)"
	test -z "$(gofmt -s -l . | tee /dev/stderr)"
