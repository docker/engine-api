// +build experimental

package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// PluginList returns the installed plugins
func (cli *Client) PluginList(ctx context.Context) (types.PluginsListResponse, error) {
	var plugins types.PluginsListResponse
	resp, err := cli.get(ctx, "/plugins", nil, nil)
	if err != nil {
		return plugins, err
	}

	err = json.NewDecoder(resp.body).Decode(&plugins)
	ensureReaderClosed(resp)
	return plugins, err
}

// PluginRemove removes a plugin
func (cli *Client) PluginRemove(ctx context.Context, name string) error {
	resp, err := cli.delete(ctx, "/plugins/"+name, nil, nil)
	ensureReaderClosed(resp)
	return err
}

// PluginEnable enables a plugin
func (cli *Client) PluginEnable(ctx context.Context, name string) error {
	resp, err := cli.post(ctx, "/plugins/"+name+"/enable", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}

// PluginDisable disables a plugin
func (cli *Client) PluginDisable(ctx context.Context, name string) error {
	resp, err := cli.post(ctx, "/plugins/"+name+"/disable", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}

// PluginInstall installs a plugin
func (cli *Client) PluginInstall(ctx context.Context, name, registryAuth string, acceptAllPermissions, noEnable bool, in io.ReadCloser, out io.Writer) error {
	headers := map[string][]string{"X-Registry-Auth": {registryAuth}}
	resp, err := cli.post(ctx, "/plugins/pull", url.Values{"name": []string{name}}, nil, headers)
	if err != nil {
		ensureReaderClosed(resp)
		return err
	}
	var privileges types.PluginPrivileges
	if err := json.NewDecoder(resp.body).Decode(&privileges); err != nil {
		return err
	}
	ensureReaderClosed(resp)

	if !acceptAllPermissions && len(privileges) > 0 {

		fmt.Fprintf(out, "Plugin %q requested the following privileges:\n", name)
		for _, privilege := range privileges {
			fmt.Fprintf(out, " - %s: %v\n", privilege.Value)
		}

		fmt.Fprint(out, "Do you grant the above permissions? [y/N] ")
		reader := bufio.NewReader(in)
		line, _, err := reader.ReadLine()
		if err != nil {
			return err
		}
		if strings.ToLower(string(line)) != "y" {
			resp, _ := cli.delete(ctx, "/plugins/"+name, nil, nil)
			ensureReaderClosed(resp)
			return pluginPermissionDenied{name}
		}
	}
	if noEnable {
		return nil
	}
	return cli.PluginEnable(ctx, name)
}

// PluginPush pushes a plugin to a registry
func (cli *Client) PluginPush(ctx context.Context, name string, registryAuth string) error {
	headers := map[string][]string{"X-Registry-Auth": {registryAuth}}
	resp, err := cli.post(ctx, "/plugins/"+name+"/push", nil, nil, headers)
	ensureReaderClosed(resp)
	return err
}

// PluginInspect inspects an existing plugin
func (cli *Client) PluginInspect(ctx context.Context, name string) (*types.Plugin, error) {
	var p types.Plugin
	resp, err := cli.get(ctx, "/plugins/"+name, nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(resp.body).Decode(&p)
	ensureReaderClosed(resp)
	return &p, err
}

// PluginSet modifies settings for an existing plugin
func (cli *Client) PluginSet(ctx context.Context, name string, args []string) error {
	resp, err := cli.post(ctx, "/plugins/"+name+"/set", nil, args, nil)
	ensureReaderClosed(resp)
	return err
}
