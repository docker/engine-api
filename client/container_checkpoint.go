package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerCheckpoint checkpoints a running container
func (cli *Client) ContainerCheckpoint(ctx context.Context, containerID string, exit bool) (types.ContainerCheckpointResponse, error) {
	query := url.Values{}
	query.Set("exit", "0")
	if exit {
		query.Set("exit", "1")
	}

	var response types.ContainerCheckpointResponse
	resp, err := cli.post(ctx, "/containers/"+containerID+"/checkpoint", query, nil, nil)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
