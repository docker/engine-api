package backend

import (
	"io"
	"net/http"
	"time"
)

// ContainersConfig is the filtering specified by the user to iterate over containers.
type ContainersConfig struct {
	All     bool   // If true show all containers, otherwise only running containers.
	Since   string // Show all containers created after this container id
	Before  string // Show all containers created before this container id
	Limit   int    // Number of containers to return at most
	Size    bool   // If true include the sizes of the containers in the response
	Filters string // Return only containers that match the filters
}

// ContainerAttachWithLogsConfig holds the streams to use when connecting to a container to view logs.
type ContainerAttachWithLogsConfig struct {
	Hijacker   http.Hijacker
	Upgrade    bool
	UseStdin   bool
	UseStdout  bool
	UseStderr  bool
	Logs       bool
	Stream     bool
	DetachKeys []byte
}

// ContainerWsAttachWithLogsConfig attach with websockets, since all
// stream data is delegated to the websocket to handle there.
type ContainerWsAttachWithLogsConfig struct {
	InStream   io.ReadCloser // Reader to attach to stdin of container
	OutStream  io.Writer     // Writer to attach to stdout of container
	ErrStream  io.Writer     // Writer to attach to stderr of container
	Logs       bool          // If true return log output
	Stream     bool          // If true return stream output
	DetachKeys []byte
}

// ContainerLogsConfig holds configs for logging operations. Exists
// for users of the daemon to to pass it a logging configuration.
type ContainerLogsConfig struct {
	Follow     bool      // If true stream log output
	Timestamps bool      // If true include timestamps for each line of log output
	Tail       string    // Return that many lines of log output from the end
	Since      time.Time // Filter logs by returning only entries after this time
	UseStdout  bool      // Whether or not to show stdout output in addition to log entries
	UseStderr  bool      // Whether or not to show stderr output in addition to log entries
	OutStream  io.Writer
	Stop       <-chan bool
}

// ContainerStatsConfig holds information for configuring the runtime
// behavior of a daemon.ContainerStats() call.
type ContainerStatsConfig struct {
	Stream    bool
	OutStream io.Writer
	Stop      <-chan bool
	Version   string
}
