package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/engine-api/client/transport"
	"github.com/docker/engine-api/types"
)

func containerCreateWithNameMock(req *http.Request) (*http.Response, error) {
	name := req.URL.Query().Get("name")
	if name != "container_name" {
		return nil, fmt.Errorf("container name not set in URL query properly. Expected `container_name`, got %s", name)
	}
	b, err := json.Marshal(types.ContainerCreateResponse{
		ID: "container_id",
	})
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(b)),
	}, nil
}

func TestContainerCreateWithName(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, containerCreateWithNameMock),
	}

	r, err := client.ContainerCreate(nil, nil, nil, "container_name")
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "container_id" {
		t.Fatalf("expected `container_id`, got %s", r.ID)
	}
}
