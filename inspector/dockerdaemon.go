package inspector

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/whalebrew/whalebrew/client"
)

// DockerDaemon fetches image inspect from docker daemon
type DockerDaemon struct {
	Client *client.Client
}

// Inspect performs requests image details from docker daemon
func (dd *DockerDaemon) Inspect(ctx context.Context, image string) (*types.ImageInspect, error) {
	return dd.Client.ImageInspect(ctx, image)
}
