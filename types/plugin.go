// +build experimental

package types

import (
	"encoding/json"
	"fmt"
)

// PluginInstallOptions holds parameters to install a plugin.
type PluginInstallOptions struct {
	Disabled             bool                 // Do not enable the plugin on install
	AcceptAllPermissions bool                 // Grant all permissions requested by the plugin
	RegistryAuth         string               // Registry credentials to be used (as base64 encoded JSON)
	PrivilegeFunc        RequestPrivilegeFunc // Function to request alternative registry credentials
	// AcceptPermissionFunc is passed the list of privileges
	// requested by the plugin, and has to decide whether
	// or not the requested privileges will be granted.
	AcceptPermissionsFunc func(PluginPrivileges) (bool, error)
}

// PluginConfig represents the values of settings which are potentially modifiable by a user.
type PluginConfig struct {
	Mounts  []PluginMount  // List of configured volumes
	Env     []string       // List of configured environment variables
	Args    []string       // List of configured command line arguments
	Devices []PluginDevice // List of configured devices
}

// Plugin represents a Docker plugin for the remote API.
type Plugin struct {
	// ID of the plugin.
	ID string `json:"Id,omitempty"`

	Name     string         // Name of the plugin
	Tag      string         // Tag of the plugin
	Active   bool           // Whether the plugin is currently enabled
	Config   PluginConfig   // Runtime configuration of the plugin
	Manifest PluginManifest // Manifest of the plugin
}

// PluginsListResponse represents a list of plugins.
type PluginsListResponse []*Plugin

const (
	authzDriver   = "AuthzDriver"
	graphDriver   = "GraphDriver"
	ipamDriver    = "IpamDriver"
	networkDriver = "NetworkDriver"
	volumeDriver  = "VolumeDriver"
)

// PluginInterfaceType represents a type that a plugin implements.
type PluginInterfaceType struct {
	Prefix     string // This is always "docker"
	Capability string // Capability that specifies the interface that is supported (e.g. "network")
	Version    string // Plugin API version. Depends on the capability.
}

// UnmarshalJSON implements json.Unmarshaler for PluginInterfaceType.
func (t *PluginInterfaceType) UnmarshalJSON(p []byte) error {
	versionIndex := len(p)
	prefixIndex := 0
	if len(p) < 2 || p[0] != '"' || p[len(p)-1] != '"' {
		return fmt.Errorf("%q is not a plugin interface type", p)
	}
	p = p[1 : len(p)-1]
loop:
	for i, b := range p {
		switch b {
		case '.':
			prefixIndex = i
		case '/':
			versionIndex = i
			break loop
		}
	}
	t.Prefix = string(p[:prefixIndex])
	t.Capability = string(p[prefixIndex+1 : versionIndex])
	if versionIndex < len(p) {
		t.Version = string(p[versionIndex+1:])
	}
	return nil
}

// MarshalJSON implements json.Marshaler for PluginInterfaceType.
func (t *PluginInterfaceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements fmt.Stringer for PluginInterfaceType
func (t PluginInterfaceType) String() string {
	return fmt.Sprintf("%s.%s/%s", t.Prefix, t.Capability, t.Version)
}

// PluginInterface describes the interface between Docker and plugin
type PluginInterface struct {
	Types  []PluginInterfaceType // Interfaces used for communication between Docker and plugin
	Socket string                // Socket used for communication between Docker and plugin
}

// PluginSetting is to be embedded in other structs,
// if they are supposed to be modifiable by the user.
type PluginSetting struct {
	Name        string   // Name of the associated configuration variable
	Description string   // Description of the associated configuration variable
	Settable    []string // List of fields that are settable by the user
}

// PluginNetwork represents the network configuration for a plugin.
type PluginNetwork struct {
	Type string
}

// PluginMount represents the mount configuration for a plugin.
type PluginMount struct {
	PluginSetting
	Source      *string  // Source path of the mount
	Destination string   // Destination path of the mount inside the container
	Type        string   // Kind of the mount (e.g. "bind")
	Options     []string // Mount options (fstab style options)
}

// PluginEnv represents an environment variable for a plugin.
type PluginEnv struct {
	PluginSetting         // Name of the environment variable is specified by PluginSettings.Name
	Value         *string // Value of the environment variable
}

// PluginArgs represents the command line arguments for a plugin.
type PluginArgs struct {
	PluginSetting
	Value []string // List of command line arguments for the plugin
}

// PluginDevice represents a host device to be mounted for a plugin.
type PluginDevice struct {
	PluginSetting
	Path *string // File path of the host device
}

// PluginUser represents the user for the plugin's process.
type PluginUser struct {
	UID uint32 `json:"Uid,omitempty"` // User ID of the plugin process
	GID uint32 `json:"Gid,omitempty"` // Group ID of the plugin process
}

// PluginManifest represents the manifest of a plugin.
type PluginManifest struct {
	ManifestVersion string          // Schema version of the manifest
	Description     string          // Description of the plugin
	Documentation   string          // Documentation for the plugins
	Interface       PluginInterface // Interfaces supported by the plugin
	Entrypoint      []string        // Default entrypoint executable and arguments
	Workdir         string          // Default working directory
	User            PluginUser      `json:",omitempty"`
	Network         PluginNetwork   // Default network configuration
	Capabilities    []string        // Capabilities supported by the plugin
	Mounts          []PluginMount   // List of default mounts
	Devices         []PluginDevice  // List of host devices to mount by default
	Env             []PluginEnv     // List of default environment variables
	Args            PluginArgs      // List of default plugin arguments
}

// PluginPrivilege describes a permission the user
// has to accept upon installing a plugin.
type PluginPrivilege struct {
	Name        string   // Name of the requested privilege (e.g. "network", "device")
	Description string   // Description of the requested privilege
	Value       []string // Further specification of the requested privilege (e.g. path of device)
}

// PluginPrivileges represents a list of plugin privileges.
type PluginPrivileges []PluginPrivilege
