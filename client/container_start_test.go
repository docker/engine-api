package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/engine-api/client/transport"
)

func TestContainerStartError(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, transport.ErrorMock(http.StatusInternalServerError, "Server error")),
	}
	err := client.ContainerStart("nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerStart(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}),
	}

	err := client.ContainerStart("container_id")
	if err != nil {
		t.Fatal(err)
	}
}
