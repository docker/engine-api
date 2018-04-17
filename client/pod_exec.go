package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hyperhq/hyper-api/types"
)

// PodExecCreate creates a new exec configuration to run an exec process.
func (cli *Client) PodExecCreate(ctx context.Context, pod, container string, config types.ExecConfig) (types.PodExecCreateResponse, error) {
	var response types.PodExecCreateResponse
	query := url.Values{}
	query["container"] = []string{container}
	resp, err := cli.post(ctx, fmt.Sprintf("/namespaces/default/pods/%v/exec", pod), query, config, nil)
	if err != nil {
		return response, err
	}
	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}

// PodExecStart starts an exec process already created in the docker host.
func (cli *Client) PodExecStart(ctx context.Context, execID string, config types.ExecStartCheck) error {
	resp, err := cli.post(ctx, fmt.Sprintf("/exec/%s/start", execID), nil, config, nil)
	ensureReaderClosed(resp)
	return err
}

// PodExecAttach attaches a connection to an exec process in the server.
// It returns a types.HijackedConnection with the hijacked connection
// and the a reader to get output. It's up to the called to close
// the hijacked connection by calling types.HijackedResponse.Close.
func (cli *Client) PodExecAttach(ctx context.Context, execID string, config types.ExecConfig) (types.HijackedResponse, error) {
	headers := map[string][]string{"Content-Type": {"application/json"}}
	return cli.postHijacked(ctx, fmt.Sprintf("/exec/%s/start", execID), nil, config, headers)
}

// PodExecInspect returns information about a specific exec process on the docker host.
func (cli *Client) PodExecInspect(ctx context.Context, execID string) (types.PodExecInspect, error) {
	var response types.PodExecInspect
	resp, err := cli.get(ctx, fmt.Sprintf("/exec/%s/json", execID), nil, nil)
	if err != nil {
		return response, err
	}

	err = json.NewDecoder(resp.body).Decode(&response)
	ensureReaderClosed(resp)
	return response, err
}
