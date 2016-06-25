package types

import (
	"io"

	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/go-units"
)

// CheckpointCreateOptions holds parameters to create a checkpoint from a container
type CheckpointCreateOptions struct {
	CheckpointID string // ID of the created checkpoint
	Exit         bool   // Whether to exit the container after the checkpoint was was successfully created
}

// ContainerAttachOptions holds parameters to attach to a container.
type ContainerAttachOptions struct {
	Stream     bool   // If true stream log output and allow input via stdin
	Stdin      bool   // Attach to standard input, enables user interaction
	Stdout     bool   // Attach to standard output
	Stderr     bool   // Attach to standard error
	DetachKeys string // Escape keys for detach
}

// ContainerCommitOptions holds parameters to commit changes of a container.
type ContainerCommitOptions struct {
	Pause     bool              // Pause container during commit
	Reference string            // Repository (and optionally tag) for the new image
	Author    string            // Author of the new image
	Comment   string            // Comment for the new image
	Changes   []string          // Dockerfile instructions to apply while committing
	Config    *container.Config // Container configuration for the new image
}

// ContainerExecInspect holds information returned by exec inspect.
type ContainerExecInspect struct {
	ExecID      string // ID of the exec instance
	ContainerID string // ID of the container
	Running     bool   // Whether the exec'ed process is running
	ExitCode    int    // Exit code of the exec'ed process
}

// ContainerListOptions holds parameters to list containers with.
type ContainerListOptions struct {
	Quiet  bool         // Only display numeric container IDs
	Size   bool         // Show container sizes in the output
	All    bool         // Also show non-running containers
	Latest bool         // Show recently created containers (including non-running ones)
	Since  string       // Only show containers created since ID (including non-running ones)
	Before string       // Only show containers created before ID (including non-running ones)
	Limit  int          // Show at most N recently created containers (including non-running ones)
	Filter filters.Args // Filter output based on these criteria
}

// ContainerLogsOptions holds parameters to filter logs with.
type ContainerLogsOptions struct {
	Follow     bool   // If true stream log output
	ShowStdout bool   // Whether or not to show stdout output in addition to log entries
	ShowStderr bool   // Whether or not to show stderr output in addition to log entries
	Timestamps bool   // If true include timestamps for each line of log output
	Details    bool   // If true include extra details provided to logs
	Tail       string // Return that many lines of log output from the end
	Since      string // Filter logs by returning only entries after this time
}

// ContainerRemoveOptions holds parameters to remove containers.
type ContainerRemoveOptions struct {
	RemoveVolumes bool // Whether to remove the volumes associated with the container
	RemoveLinks   bool // Whether to remove the link with the specified name, instead of removing a container with that name
	Force         bool // Whether to kill the container if it is running, instead of not removing it
}

// ContainerStartOptions holds parameters to start containers.
type ContainerStartOptions struct {
	CheckpointID string // Checkpoint to start the container from
}

// CopyToContainerOptions holds information
// about files to copy into a container
type CopyToContainerOptions struct {
	AllowOverwriteDirWithFile bool // Allow overwriting an existing directory with a non-directory and vice versa
}

// EventsOptions hold parameters to filter events with.
type EventsOptions struct {
	Since   string       // Only show events created since this timestamp
	Until   string       // Stream events until this timestamp
	Filters filters.Args // Filter output based on these criteria
}

// NetworkListOptions holds parameters to filter the list of networks with.
type NetworkListOptions struct {
	Filters filters.Args // Filter output based on these criteria
}

// ImageBuildOptions holds the information
// necessary to build images.
type ImageBuildOptions struct {
	Tags           []string
	SuppressOutput bool
	RemoteContext  string
	NoCache        bool
	Remove         bool
	ForceRemove    bool
	PullParent     bool
	Isolation      container.Isolation
	CPUSetCPUs     string
	CPUSetMems     string
	CPUShares      int64
	CPUQuota       int64
	CPUPeriod      int64
	Memory         int64
	MemorySwap     int64
	CgroupParent   string
	ShmSize        int64
	Dockerfile     string
	Ulimits        []*units.Ulimit
	BuildArgs      map[string]string
	AuthConfigs    map[string]AuthConfig
	Context        io.Reader
	Labels         map[string]string
}

// ImageBuildResponse holds information returned by a server after building an image.
type ImageBuildResponse struct {
	Body   io.ReadCloser // Body must be closed to avoid a resource leak
	OSType string        // Operating system type of the Docker daemon
}

// ImageCreateOptions holds information to create images.
type ImageCreateOptions struct {
	RegistryAuth string // Registry credentials to be used (as base64 encoded JSON)
}

// ImageImportSource holds source information for ImageImport.
type ImageImportSource struct {
	Source     io.Reader // Data sent to the server to create the image from (mutually exclusive with SourceName)
	SourceName string    // Name of the image to pull (mutually exclusive with Source)
}

// ImageImportOptions holds information to import images from the client host.
type ImageImportOptions struct {
	Tag     string   // Tag is the name to tag this image with. This attribute is deprecated.
	Message string   // Message is the message to tag the image with
	Changes []string // Dockerfile instructions to apply while importing
}

// ImageListOptions holds parameters to filter the list of images with.
type ImageListOptions struct {
	MatchName string       // Only return images with matching names
	All       bool         // Whether to include intermediate (i.e. untagged) images in the output
	Filters   filters.Args // Filter output based on these criteria
}

// ImageLoadResponse returns information to the client about a load process.
type ImageLoadResponse struct {
	Body io.ReadCloser // Body must be closed to avoid a resource leak
	JSON bool          // Whether the response body contains JSON instead of plain text
}

// ImagePullOptions holds information to pull images.
type ImagePullOptions struct {
	// Pull all tags from the specified repository,
	// even if a repository:tag combination is specified
	All bool

	RegistryAuth  string               // Registry credentials to be used (as base64 encoded JSON)
	PrivilegeFunc RequestPrivilegeFunc // Function to request alternative registry credentials
}

// ImagePushOptions holds information to push images.
type ImagePushOptions ImagePullOptions

// RequestPrivilegeFunc is a function interface that clients can
// supply to retry operations after getting an authorization error.
//
// This function returns the new registry authentication header value
// in Base64 format, or an error if the privilege request fails.
type RequestPrivilegeFunc func() (string, error)

// ImageRemoveOptions holds parameters to remove images.
type ImageRemoveOptions struct {
	Force         bool // Also remove all tags referencing the image, instead of failing when such tags exist
	PruneChildren bool // Also remove all untagged parents of the image
}

// ImageSearchOptions holds parameters to search images with.
type ImageSearchOptions struct {
	RegistryAuth  string               // Registry credentials to be used (as base64 encoded JSON)
	PrivilegeFunc RequestPrivilegeFunc // Function to request alternative registry credentials
	Limit         int                  // Maximum number of search results
	Filters       filters.Args         // Filter output based on these criteria
}

// ResizeOptions holds parameters to resize a TTY.
// It can be used to resize container TTYs and exec process TTYs too.
type ResizeOptions struct {
	Height int
	Width  int
}

// VersionResponse holds version information for the client and the server.
type VersionResponse struct {
	Client *Version
	Server *Version
}

// ServerOK returns true when the client could connect to the docker server
// and parse the information received. It returns false otherwise.
func (v VersionResponse) ServerOK() bool {
	return v.Server != nil
}

// NodeListOptions holds parameters to list nodes with.
type NodeListOptions struct {
	Filter filters.Args // Filter output based on these criteria
}

// ServiceCreateResponse contains the information returned to a client
// on the  creation of a new service.
type ServiceCreateResponse struct {
	ID string // ID of the created service
}

// ServiceListOptions holds parameters to list services with.
type ServiceListOptions struct {
	Filter filters.Args // Filter output based on these criteria
}

// TaskListOptions holds parameters to list tasks with.
type TaskListOptions struct {
	Filter filters.Args // Filter output based on these criteria
}
