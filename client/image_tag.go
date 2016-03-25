package client

import (
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ImageTag tags an image in the docker host
func (cli *Client) ImageTag(ctx context.Context, imageID, ref string, options types.ImageTagOptions) error {
	repository, tag, err := parseReference(ref)
	if err != nil {
		return err
	}

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
