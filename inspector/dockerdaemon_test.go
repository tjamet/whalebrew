package inspector_test

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/whalebrew/whalebrew/client"
	"github.com/whalebrew/whalebrew/inspector"
)

func TestDockerDaemon(t *testing.T) {
	c, err := client.NewClient()
	assert.NoError(t, err)
	r := inspector.DockerDaemon{
		Client: c,
	}
	inspect, err := r.Inspect(context.Background(), "whalebrew/awscli:latest")
	assert.NoError(t, err)
	assert.NotNil(t, inspect)
	assert.NotNil(t, inspect.Config)
	assert.Equal(t, "aws", inspect.Config.Labels["io.whalebrew.name"])
}
