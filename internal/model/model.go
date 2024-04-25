package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type NodeResourceUsage struct {
	RamUsage          float64 `json:"ram_usage"`
	CpuUsage          float64 `json:"cpu_usage"`
	RamTotal          float64 `json:"ram_total"`
	CpuTotal          float64 `json:"cpu_total"`
	RunningContainers int     `json:"running_containers"`
}

type BasicContainerStatistics struct {
	Name          string              `json:"name"`
	Id            string              `json:"id"`
	Image         string              `json:"image"`
	ResourceUsage *BasicResourceUsage `json:"resource_usage"`
	RamTotal      float64             `json:"ram_total"`
	CpuTotal      float64             `json:"cpu_total"`
	//Port       int     `json:"port"`
	StartCmd  string `json:"start_cmd"`
	IsRunning bool   `json:"is_running"`
}

type BasicResourceUsage struct {
	//DiskUsage  float64 `json:"disk_usage"`
	RamUsage float64 `json:"ram_usage"`
	CpuUsage float64 `json:"cpu_usage"`
}

type CreateContainerRequest struct {
	Name    string
	Config  container.Config
	Host    container.HostConfig
	Network network.NetworkingConfig
}

// Docker image utilized, the image startup command, the worker’s name, allocated ports,
// and resource limits, including disk usage, RAM usage, and CPU usage for each individual worker node.

type CreateSimpleContainerRequest struct {
	Image         string            `json:"image"`
	Cmd           []string          `json:"cmd"`
	Name          string            `json:"name"`
	Host          *SimpleHostConfig `json:"config"`
	NeedsGpu      bool              `json:"needs_gpu"`
	StartOnCreate bool              `json:"start_on_create"`
}

type SimpleHostConfig struct {
	PortBindings map[int][]int        `json:"port_bindings"`
	Resources    *container.Resources `json:"resources"`
	AutoRemove   bool                 `json:"auto_remove"`
}

type Container struct {
	Id        string
	Container types.Container
}
