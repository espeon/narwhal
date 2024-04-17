package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

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
