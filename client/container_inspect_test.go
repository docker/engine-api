package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/engine-api/client/transport"
	"github.com/docker/engine-api/types"
)

func TestContainerInspectError(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, transport.ErrorMock(http.StatusInternalServerError, "Server error")),
	}

	_, err := client.ContainerInspect("nothing")
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}

func TestContainerInspectContainerNotFound(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, transport.ErrorMock(http.StatusNotFound, "Server error")),
	}

	_, err := client.ContainerInspect("unknown")
	if err == nil || !IsErrContainerNotFound(err) {
		t.Fatalf("expected a containerNotFound error, got %v", err)
	}
}

func TestContainerInspect(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, func(req *http.Request) (*http.Response, error) {
			content, err := json.Marshal(types.ContainerJSON{
				ContainerJSONBase: &types.ContainerJSONBase{
					ID:    "container_id",
					Image: "image",
					Name:  "name",
				},
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(content)),
			}, nil
		}),
	}

	r, err := client.ContainerInspect("container_id")
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "container_id" {
		t.Fatalf("expected `container_id`, got %s", r.ID)
	}
	if r.Image != "image" {
		t.Fatalf("expected `image`, got %s", r.ID)
	}
	if r.Name != "name" {
		t.Fatalf("expected `name`, got %s", r.ID)
	}
}
