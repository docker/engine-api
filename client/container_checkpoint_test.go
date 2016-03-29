package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func TestContainerCheckpointError(t *testing.T) {
	client := &Client{
		transport: newMockClient(nil, errorMock(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.ContainerCheckpoint(context.Background(), "nothing", true)

	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerCheckpoint(t *testing.T) {
	expectedContainerID := "container_id"
	expectedCheckpointID := "checkpoint_id"

	client := &Client{
		transport: newMockClient(nil, func(req *http.Request) (*http.Response, error) {
			exit := req.URL.Query().Get("exit")
			if exit != "1" {
				return nil, fmt.Errorf("exit not set in URL query properly. Expected '1', got %s", exit)
			}
			b, err := json.Marshal(types.ContainerCheckpointResponse{
				ID: expectedCheckpointID,
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			}, nil
		}),
	}

	r, err := client.ContainerCheckpoint(context.Background(), expectedContainerID, true)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != expectedCheckpointID {
		t.Fatalf("expected `checkpoint_id`, got %s", r.ID)
	}
}
