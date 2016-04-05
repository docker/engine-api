# Docker API swagger spec

This swagger spec was created by hand and based on Docker's public remote API documentation.  Several open source projects contributed to it's creation.

A note about swagger.  In case you do not know what Swagger is, it is an open source tool used to specify how clients and servers talk to each other via REST/JSON.  It provide a way to generate both a standard REST server and a client component that other software may use to talk to the REST server.  For more information on Swagger and go-swagger, please visit the following sites.

- [swagger](http://swagger.io/)
- [go-swagger](https://goswagger.io/)

For this spec, go-swagger was used to generate golang code.  As much attempt to make this spec and the generated server comply with the existing engine-api server behavior, it is not perfect.  The server created using this spec will result in a clean, pure REST server.  It cannot cover all the edge cases and specific behaviors of the current engine-api server.  Over time, the engine-api server can migrate over to the swagger-generated server or this spec can just be used to generate clients for various languages.

## Instructions

To build, use the top level makefile.

For the server:
```
$ make deps
$ make dockerapi-server
```
NOTE:  Currently, the make clean target cleans up all files created during a server generation.  Some files shouldn't be removed, such as configure_docker.go.  This will soon be fixed, but it is recommended you copy the server code out for now.

For the client:
```
$ make deps
$ make dockerapi-client
```