package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"
)

// ServiceInspectWithRaw returns the service information and the raw data.
func (cli *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string) (swarm.InspectService, []byte, error) {
	serverResp, err := cli.get(ctx, "/services/"+serviceID, nil, nil)
	if err != nil {
		if serverResp.statusCode == http.StatusNotFound {
			return swarm.InspectService{}, nil, serviceNotFoundError{serviceID}
		}
		return swarm.InspectService{}, nil, err
	}
	defer ensureReaderClosed(serverResp)

	body, err := ioutil.ReadAll(serverResp.body)
	if err != nil {
		return swarm.InspectService{}, nil, err
	}

	var response swarm.InspectService
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&response)
	return response, body, err
}
