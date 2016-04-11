package client

import (
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerRestore restores a running container
func (cli *Client) ContainerRestore(ctx context.Context, options types.ContainerRestoreOptions) error {
	query := url.Values{}
	query.Set("id", options.CheckpointID)

	_, err := cli.post(ctx, "/containers/"+options.ContainerID+"/restore", query, nil, nil)

	return err
}
