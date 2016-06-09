package swarm

import "time"

// ContainerSpec represents the spec of a container.
type ContainerSpec struct {
	Image           string            `json:",omitempty"`
	Labels          map[string]string `json:",omitempty"`
	Command         []string          `json:",omitempty"`
	Args            []string          `json:",omitempty"`
	Env             []string          `json:",omitempty"`
	Dir             string            `json:",omitempty"`
	User            string            `json:",omitempty"`
	Mounts          []Mount           `json:",omitempty"`
	StopGracePeriod *time.Duration    `json:",omitempty"`
}

// MountType represents the type of a mount.
type MountType string

const (
	// MountTypeBind BIND
	MountTypeBind MountType = "BIND"
	// MountTypeVolume VOLUME
	MountTypeVolume MountType = "VOLUME"
)

// Mount represents a mount (volume).
type Mount struct {
	Type     MountType `json:",omitempty"`
	Source   string    `json:",omitempty"`
	Target   string    `json:",omitempty"`
	Writable bool      `json:",omitempty"`

	BindOptions   *BindOptions   `json:",omitempty"`
	VolumeOptions *VolumeOptions `json:",omitempty"`
}

// MountPropagation represents the propagation of a mount.
type MountPropagation string

const (
	// MountPropagationRPrivate RPRIVATE
	MountPropagationRPrivate MountPropagation = "RPRIVATE"
	// MountPropagationPrivate PRIVATE
	MountPropagationPrivate MountPropagation = "PRIVATE"
	// MountPropagationRShared RSHARED
	MountPropagationRShared MountPropagation = "RSHARED"
	// MountPropagationShared SHARED
	MountPropagationShared MountPropagation = "SHARED"
	// MountPropagationRSlave RSLAVE
	MountPropagationRSlave MountPropagation = "RSLAVE"
	// MountPropagationSlave SLAVE
	MountPropagationSlave MountPropagation = "SLAVE"
)

type BindOptions struct {
	Propagation MountPropagation `json:",omitempty"`
}

// VolumeOptions represents the options for a mount of type volume.
type VolumeOptions struct {
	Populate     bool              `json:",omitempty"`
	Labels       map[string]string `json:",omitempty"`
	DriverConfig Driver            `json:",omitempty"`
}
