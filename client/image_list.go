package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
	"golang.org/x/net/context"
)

// ImageList returns a list of images in the docker host.
func (cli *Client) ImageList(ctx context.Context, options types.ImageListOptions) (images []types.Image, err error) {
	query := url.Values{}

	if options.Filters.Len() > 0 {
		filterJSON, err0 := filters.ToParamWithVersion(cli.version, options.Filters)
		if err0 != nil {
			err = err0
			return
		}
		query.Set("filters", filterJSON)
	}
	if options.MatchName != "" {
		// FIXME rename this parameter, to not be confused with the filters flag
		query.Set("filter", options.MatchName)
	}
	if options.All {
		query.Set("all", "1")
	}

	serverResp, err := cli.get(ctx, "/images/json", query, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(serverResp.body).Decode(&images)
	ensureReaderClosed(serverResp)
	return
}
