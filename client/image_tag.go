package client

import (
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ImageTag tags an image in the docker host
func (cli *Client) ImageTag(ctx context.Context, imageID, repository, tag string, options types.ImageTagOptions) error {
	query := url.Values{}
	query.Set("repo", repository)
	query.Set("tag", tag)
	if options.Force {
		query.Set("force", "1")
	}

	resp, err := cli.post(ctx, "/images/"+imageID+"/tag", query, nil, nil)
	ensureReaderClosed(resp)
	return err
}
