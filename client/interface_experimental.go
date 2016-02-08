// +build experimental

package client

import (
	"github.com/docker/engine-api/types"
)

// ExperimentalAPIClient is an interface that implements Experimental API methods
type ExperimentalAPIClient interface {
	ContainerCheckpoint(containerID string, options types.CriuConfig) error
	ContainerRestore(containerID string, options types.CriuConfig, forceRestore bool) error
}
