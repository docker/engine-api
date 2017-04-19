package client

import (
	"net/url"
	"net/http"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ContainerRemove kills and removes a container from the docker host.
func (cli *Client) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	query := url.Values{}
	if options.RemoveVolumes {
		query.Set("v", "1")
	}
	if options.RemoveLinks {
		query.Set("link", "1")
	}

	if options.Force {
		query.Set("force", "1")
	}

	resp, err := cli.delete(ctx, "/containers/"+containerID, query, nil)
	if err != nil {
		if resp.statusCode == http.StatusNotFound {
			err = containerNotFoundError{containerID}
		}
	}
	ensureReaderClosed(resp)
	return err
}
