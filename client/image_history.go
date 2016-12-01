package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// ImageHistory returns the changes in an image in history format.
func (cli *Client) ImageHistory(ctx context.Context, imageID string) (history []types.ImageHistory, err error) {
	serverResp, err := cli.get(ctx, "/images/"+imageID+"/history", url.Values{}, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(serverResp.body).Decode(&history)
	ensureReaderClosed(serverResp)
	return
}
