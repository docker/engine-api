// Package types is used for API stability in the types and response to the
// consumers of the API stats endpoint.
package types

import "time"

// ThrottlingData stores the CPU throttling stats of one running container.
type ThrottlingData struct {
	Periods          uint64 `json:"periods"`           // Number of periods with throttling active
	ThrottledPeriods uint64 `json:"throttled_periods"` // Number of periods where the container hit its throttling limit
	ThrottledTime    uint64 `json:"throttled_time"`    // Aggregate time the container was throttled for in nanoseconds
}

// CPUUsage stores all CPU stats aggregated since the container was created.
type CPUUsage struct {
	TotalUsage        uint64   `json:"total_usage"`         // Total CPU time consumed (in nanoseconds)
	PercpuUsage       []uint64 `json:"percpu_usage"`        // Total CPU time consumed per core (in nanoseconds)
	UsageInKernelmode uint64   `json:"usage_in_kernelmode"` // Time spent by tasks of the cgroup in kernel mode (in nanoseconds)
	UsageInUsermode   uint64   `json:"usage_in_usermode"`   // Time spent by tasks of the cgroup in user mode (in seconds)
}

// CPUStats aggregates and wraps all CPU related info of a container.
type CPUStats struct {
	CPUUsage       CPUUsage       `json:"cpu_usage"`
	SystemUsage    uint64         `json:"system_cpu_usage"`
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}

// MemoryStats aggregates all memory stats since the container was created.
type MemoryStats struct {
	Usage    uint64 `json:"usage"`     // Current memory usage of the container
	MaxUsage uint64 `json:"max_usage"` // Maximum memory usage of the container since it was created
	// TODO(vishh): Export these as stronger types.
	Stats   map[string]uint64 `json:"stats"`   // All stats exported via memory.stat
	Failcnt uint64            `json:"failcnt"` // Number of times when memory usage hit limits
	Limit   uint64            `json:"limit"`   // Maximum allowed memory usage configured for the container
}

// BlkioStatEntry is one small entity to store a piece of Blkio stats
// TODO Windows: This can be factored out
type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}

// BlkioStats stores all IO service stats for data read and write.
// TODO Windows: This can be factored out
type BlkioStats struct {
	// Number of bytes transferred to and from the block device
	IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
	IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
	IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
	IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
	IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
	IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
	IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
	SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
}

// NetworkStats aggregates all network stats of a container.
// TODO Windows: This will require refactoring
type NetworkStats struct {
	RxBytes   uint64 `json:"rx_bytes"`
	RxPackets uint64 `json:"rx_packets"`
	RxErrors  uint64 `json:"rx_errors"`
	RxDropped uint64 `json:"rx_dropped"`
	TxBytes   uint64 `json:"tx_bytes"`
	TxPackets uint64 `json:"tx_packets"`
	TxErrors  uint64 `json:"tx_errors"`
	TxDropped uint64 `json:"tx_dropped"`
}

// PidsStats contains the stats of a container's pids.
type PidsStats struct {
	// Current is the number of pids in the cgroup.
	Current uint64 `json:"current,omitempty"`

	// Limit is the hard limit on the number of pids in the cgroup.
	// A "Limit" of 0 means that there is no limit.
	Limit uint64 `json:"limit,omitempty"`
}

// Stats is the ultimate struct aggregating all types of stats of one container.
type Stats struct {
	Read        time.Time   `json:"read"`
	PreCPUStats CPUStats    `json:"precpu_stats,omitempty"`
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`
	BlkioStats  BlkioStats  `json:"blkio_stats,omitempty"`
	PidsStats   PidsStats   `json:"pids_stats,omitempty"`
}

// StatsJSON is newly used Networks
type StatsJSON struct {
	Stats

	// Networks request version >=1.21
	Networks map[string]NetworkStats `json:"networks,omitempty"`
}
