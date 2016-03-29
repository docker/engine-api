package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

func TestContainerRestoreError(t *testing.T) {
	client := &Client{
		transport: newMockClient(nil, errorMock(http.StatusInternalServerError, "Server error")),
	}
	err := client.ContainerRestore(context.Background(), "nothing", "nothing")

	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerRestore(t *testing.T) {
	containerID := "container_id"
	checkpointID := "checkpoint_id"

	client := &Client{
		transport: newMockClient(nil, func(req *http.Request) (*http.Response, error) {
			id := req.URL.Query().Get("id")
			if id != checkpointID {
				return nil, fmt.Errorf("id not set in URL query properly. Expected 'checkpoint_id', got %s", id)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			}, nil
		}),
	}

	err := client.ContainerRestore(context.Background(), containerID, checkpointID)
	if err != nil {
		t.Fatal(err)
	}
}
