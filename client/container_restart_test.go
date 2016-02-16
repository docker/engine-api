package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/engine-api/client/transport"
)

func TestContainerRestartError(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, transport.ErrorMock(http.StatusInternalServerError, "Server error")),
	}
	err := client.ContainerRestart("nothing", 0)
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerRestart(t *testing.T) {
	expectedURL := "/containers/container_id/restart"
	client := &Client{
		transport: transport.NewMockClient(nil, func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			t := req.URL.Query().Get("t")
			if t != "100" {
				return nil, fmt.Errorf("t (timeout) not set in URL query properly. Expected '100', got %s", t)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}),
	}

	err := client.ContainerRestart("container_id", 100)
	if err != nil {
		t.Fatal(err)
	}
}
