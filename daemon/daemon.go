// The daemon package controls all interaction with the Docker daemon. Any
// options for containers, images, networks, and volumes are all routed
// through this package.

package daemon

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// DockerInterface is our primary source of information about
// the Docker daemon and associated containers, images, networks,
// and volumes.
type DockerInterface struct {
	Client     *client.Client
	Containers []types.Container
	Images     []types.ImageSummary
	Info       types.Info
	Networks   []types.NetworkResource
	Volumes    []*types.Volume
}

// RefreshContainers updates the DockerInterface's Containers fields
// with the latest information from the Docker API.
func (di *DockerInterface) RefreshContainers(ctx context.Context) error {
	var err error

	if di.Containers, err = di.Client.ContainerList(
		ctx, types.ContainerListOptions{All: true}); err != nil {
		return fmt.Errorf("failed to refresh container list: %s", err)
	}
	return nil
}

// RefreshImages updates the DockerInterface's Images field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshImages(ctx context.Context) error {
	var err error

	if di.Images, err = di.Client.ImageList(
		ctx, types.ImageListOptions{All: true}); err != nil {
		return fmt.Errorf("failed to refresh image list: %s", err)
	}
	return nil
}

// RefreshInfo updates the DockerInterface's Info field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshInfo(ctx context.Context) error {
	var err error

	if di.Info, err = di.Client.Info(ctx); err != nil {
		return fmt.Errorf("failed to refresh docker data: %s", err)
	}
	return nil
}

// RefreshNetworks updates the DockerInterface's Networks field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshNetworks(ctx context.Context) error {
	var err error

	if di.Networks, err = di.Client.NetworkList(
		ctx, types.NetworkListOptions{}); err != nil {
		return fmt.Errorf("failed to refresh network list: %s", err)
	}
	return nil
}

// RefreshVolumes updates the DockerInterface's Volumes field with the
// latest information from the Docker API.
func (di *DockerInterface) RefreshVolumes(ctx context.Context) error {
	var err error
	var volumeBody volume.VolumeListOKBody

	if volumeBody, err = di.Client.VolumeList(ctx, filters.Args{}); err != nil {
		return fmt.Errorf("failed to refresh volume list: %s", err)
	}

	di.Volumes = volumeBody.Volumes
	return nil
}

// NewInterface returns a DockerInterface with information about the
// Docker daemon and all containers, images, networks, and volumes.
func NewInterface(ctx context.Context) (*DockerInterface, error) {
	var err error
	di := &DockerInterface{}

	di.Client, err = client.NewClientWithOpts(
		client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	if err = di.RefreshContainers(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshImages(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshInfo(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshNetworks(ctx); err != nil {
		return nil, err
	}

	if err = di.RefreshVolumes(ctx); err != nil {
		return nil, err
	}
	return di, nil
}
