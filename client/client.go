package client

import (
	"context"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DefaultVersion is the Engine API version used by Whalebrew
const DefaultVersion string = "1.20"

type Client struct {
	*client.Client
}

// NewClient returns a Docker client configured for Whalebrew
func NewClient() (*Client, error) {
	dockerClient, err := client.NewClientWithOpts(client.WithVersion(DefaultVersion), client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{
		dockerClient,
	}, nil
}

func (c *Client) PullImageIfNotExists(ctx context.Context, imageName string) error {
	_, _, err := c.ImageInspectWithRaw(ctx, imageName)
	if client.IsErrNotFound(err) {
		if err = pullImage(imageName); err != nil {
			return err
		}
	}
	return err
}

func (c *Client) ImageInspect(ctx context.Context, imageName string) (*types.ImageInspect, error) {
	imageInspect, _, err := c.ImageInspectWithRaw(ctx, imageName)
	return &imageInspect, err
}

func pullImage(image string) error {
	c := exec.Command("docker", "pull", image)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
