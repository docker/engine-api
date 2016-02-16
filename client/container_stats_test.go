package client

import (
	"net/http"
	"testing"

	"github.com/docker/engine-api/client/transport"

	"golang.org/x/net/context"
)

func TestContainerStatsError(t *testing.T) {
	client := &Client{
		transport: transport.NewMockClient(nil, transport.ErrorMock(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.ContainerStats(context.Background(), "nothing", false)
	if err == nil || err.Error() != "Error response from daemon: Server error" {
		t.Fatalf("expected a Server Error, got %v", err)
	}
}
