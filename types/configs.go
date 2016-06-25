package types

import (
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
)

// configs holds structs used for internal communication between the
// frontend (such as an http server) and the backend (such as the
// docker daemon).

// ContainerCreateConfig is the parameter set to ContainerCreate()
type ContainerCreateConfig struct {
	Name             string
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	AdjustCPUShares  bool
}

// ContainerRmConfig holds arguments for the container remove operation.
// This struct is used to tell the backend what operations to perform.
type ContainerRmConfig struct {
	ForceRemove, RemoveVolume, RemoveLink bool
}

// ContainerCommitConfig contains build configs for commit operation,
// and is used when making a commit with the current state of the container.
type ContainerCommitConfig struct {
	Pause        bool              // Pause container during commit
	Repo         string            // Repository name of the new image
	Tag          string            // Tag name of the new image
	Author       string            // Author of the new image
	Comment      string            // Comment for the new image
	MergeConfigs bool              // Merge container config into commit config before commit
	Config       *container.Config // Container configuration for the new image
}

// ExecConfig is a small subset of the Config struct that
// holds the configuration for the exec feature of docker.
type ExecConfig struct {
	User         string   // User that will run the command
	Privileged   bool     // Run the command in privileged mode
	Tty          bool     // Attach standard streams to a tty
	AttachStdin  bool     // Attach the standard input, enables user interaction
	AttachStdout bool     // Attach the standard output
	AttachStderr bool     // Attach the standard error
	Detach       bool     // Execute in detached mode
	DetachKeys   string   // Escape keys for detach
	Cmd          []string // Execution commands and args
}
