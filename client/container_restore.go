package client

import (
	"net/url"

	"golang.org/x/net/context"
)

// ContainerCheckpoint checkpoints a running container
func (cli *Client) ContainerRestore(ctx context.Context, containerID string, checkpointID string) error {
	query := url.Values{}
	query.Set("id", checkpointID)

	_, err := cli.post(ctx, "/containers/"+containerID+"/restore", query, nil, nil)

	return err
}
