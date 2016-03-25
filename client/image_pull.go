package client

import (
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"github.com/docker/engine-api/types"
)

// ImagePull requests the docker host to pull an image from a remote registry.
// It executes the privileged function if the operation is unauthorized
// and it tries one more time.
// It's up to the caller to handle the io.ReadCloser and close it properly.
//
// FIXME(vdemeester): there is currently few way to use this from docker/docker
// - if not in trusted content, ref is used to pass the whole reference, and tag is empty
// - if in trusted content, ref is used to pass the reference name, and tag for the digest
func (cli *Client) ImagePull(ctx context.Context, ref, tag string, options types.ImagePullOptions) (io.ReadCloser, error) {
	query := url.Values{}
	query.Set("fromImage", ref)
	if tag != "" {
		query.Set("tag", tag)
	}

	resp, err := cli.tryImageCreate(ctx, query, options.RegistryAuth)
	if resp.statusCode == http.StatusUnauthorized {
		newAuthHeader, privilegeErr := options.PrivilegeFunc()
		if privilegeErr != nil {
			return nil, privilegeErr
		}
		resp, err = cli.tryImageCreate(ctx, query, newAuthHeader)
	}
	if err != nil {
		return nil, err
	}
	return resp.body, nil
}
