package client

import (
	"io"
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ImageImport creates a new image based in the source options.
// It returns the JSON content in the response body.
func (cli *Client) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	repository, tag, err := parseReference(ref)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("fromSrc", source.SourceName)
	query.Set("repo", repository)
	query.Set("tag", tag)
	query.Set("message", options.Message)
	for _, change := range options.Changes {
		query.Add("changes", change)
	}

	resp, err := cli.postRaw(ctx, "/images/create", query, source.Source, nil)
	if err != nil {
		return nil, err
	}
	return resp.body, nil
}
