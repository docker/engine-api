.PHONY: all deps test validate

all: deps test validate

SWAGGER ?= $(GOPATH)/bin/swagger$(BIN_ARCH)

swagger: $(SWAGGER)

deps:
	go get -t ./...
	go get github.com/golang/lint/golint

test:
	go test -race -cover ./...

validate:
	go vet ./...
	test -z "$(golint ./... | tee /dev/stderr)"
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

$(SWAGGER):
	@echo building $(SWAGGER)...
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	
dockerapi-server: swagger/swagger.json $(SWAGGER)
	@echo regenerating swagger models and operations for Docker API server...
	@$(SWAGGER) generate server -A docker -t $(dir $<)/server -f $< 
	go build -o ./binary/docker-server ./swagger/server/cmd/docker-server

dockerapi-client: swagger/swagger.json $(SWAGGER)
	@echo regenerating swagger client for Docker API server...
	@$(SWAGGER) generate client -A docker -t $(dir $<)/client -f $<
	go build -o ./binary/docker-client ./swagger/client/cmd/docker

clean:
	rm -rf ./binary

	@echo removing swagger generated files...
	rm -rf ./swagger/server
	rm -rf ./swagger/client
