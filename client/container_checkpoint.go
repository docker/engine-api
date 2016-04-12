package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerCheckpoint checkpoints a running container
func (cli *Client) ContainerCheckpoint(ctx context.Context, options types.ContainerCheckpointOptions) (types.ContainerCheckpointResponse, error) {
	query := url.Values{}
	query.Set("id", options.CheckpointID)
	query.Set("exit", "0")
	if options.Exit {
		query.Set("exit", "1")
	}

	var response types.ContainerCheckpointResponse
	resp, err := cli.post(ctx, "/containers/"+options.ContainerID+"/checkpoint", query, nil, nil)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
