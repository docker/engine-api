package client

import (
	"context"
	"fmt"

	"github.com/hyperhq/hyper-api/types"
)

// PodExecResize changes the size of the tty for an exec process running inside a container.
func (cli *Client) PodExecResize(ctx context.Context, execID string, options types.ResizeOptions) error {
	return cli.resize(ctx, fmt.Sprintf("/exec/%s", execID), options.Height, options.Width)
}
