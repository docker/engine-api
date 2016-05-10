package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func TestInfo(t *testing.T) {
	expectedURL := "/info"
	client := &Client{
		transport: newMockClient(nil, func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}
			info := &types.Info{
				ID:         "daemonID",
				Containers: 3,
			}
			b, err := json.Marshal(info)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(b)),
			}, nil
		}),
	}

	info, err := client.Info(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if info.ID != "daemonID" {
		t.Fatalf("expected daemonID, got %s", info.ID)
	}

	if info.Containers != 3 {
		t.Fatalf("expected 3 containers, got %d", info.Containers)
	}
}
