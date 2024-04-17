package service

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"nat.vg/narwhal/internal/model"
)

type Service interface {
	List(ctx context.Context) ([]types.Container, error)
	Get(ctx context.Context, id string) (types.ContainerJSON, error)
	Create(ctx context.Context, req model.CreateContainerRequest) error
	CreateSimple(ctx context.Context, req model.CreateSimpleContainerRequest) (container.CreateResponse, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Remove(ctx context.Context, id string, force bool, removeVolumes bool) error
}

type DockerService struct {
	cli *client.Client
}

func NewNarwhalService() Service {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return DockerService{
		cli: cli,
	}
}

// List containers service
func (s DockerService) List(ctx context.Context) ([]types.Container, error) {
	containers, err := s.cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	return containers, nil
}

// get container
func (s DockerService) Get(ctx context.Context, id string) (types.ContainerJSON, error) {
	container, err := s.cli.ContainerInspect(ctx, id)
	if err != nil {
		panic(err)
	}
	return container, nil
}

func createPortBindings(ports map[int][]int) nat.PortMap {
	portMap := make(nat.PortMap)

	for hostPort, containerPorts := range ports {
		for _, containerPort := range containerPorts {
			port := nat.Port(fmt.Sprintf("%d/tcp", containerPort))
			portMap[port] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", hostPort)}}
		}
	}

	return portMap
}

func (s DockerService) CreateSimple(ctx context.Context, req model.CreateSimpleContainerRequest) (container.CreateResponse, error) {
	// build our config
	config := container.Config{
		Image:    req.Image,
		Cmd:      req.Cmd,
		Hostname: req.Name,
	}

	bindings := createPortBindings(req.Host.PortBindings)
	println(bindings["8080/tcp"])
	hostConfig := container.HostConfig{
		PortBindings: bindings,
		Resources:    *req.Host.Resources,
		AutoRemove:   req.Host.AutoRemove,
	}
	println("ball")

	// pull image if necessary
	if _, _, err := s.cli.ImageInspectWithRaw(ctx, req.Image); err != nil {
		println("pulling image")
		reader, err := s.cli.ImagePull(ctx, req.Image, types.ImagePullOptions{})
		if err != nil {
			// report error
			fmt.Fprintln(os.Stderr, "Error pulling image:", err.Error())
			return container.CreateResponse{}, fmt.Errorf("%s", err.Error())
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)

		println("pulled image")
	}

	res, err := s.cli.ContainerCreate(ctx, &config, &hostConfig, nil, nil, req.Name)
	if err != nil {
		// report error
		fmt.Fprintln(os.Stderr, "Error creating container:", err.Error())
		return res, fmt.Errorf("%s", err.Error())
	}
	// autostart container
	if req.StartOnCreate {
		err = s.cli.ContainerStart(ctx, res.ID, container.StartOptions{})
		if err != nil {
			// report error
			fmt.Fprintln(os.Stderr, "Error starting container:", err.Error())
			return res, fmt.Errorf("%s", err.Error())
		}
	}
	return res, nil
}

// create container
func (s DockerService) Create(ctx context.Context, req model.CreateContainerRequest) error {
	_, err := s.cli.ContainerCreate(ctx, &req.Config, &req.Host, &req.Network, nil, req.Name)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s DockerService) Start(ctx context.Context, id string) error {
	err := s.cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}

func (s DockerService) Stop(ctx context.Context, id string) error {
	stopOpts := container.StopOptions{
		// 10 seconds
		Timeout: &[]int{10}[0],
	} // TODO: get this from config
	err := s.cli.ContainerStop(ctx, id, stopOpts)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}

func (s DockerService) Remove(ctx context.Context, id string, force bool, removeVolumes bool) error {
	removeOpts := container.RemoveOptions{
		RemoveVolumes: removeVolumes,
		Force:         force,
	}
	err := s.cli.ContainerRemove(ctx, id, removeOpts)
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	return nil
}
