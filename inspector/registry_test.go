package inspector_test

import (
	"testing"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/whalebrew/whalebrew/inspector"
)

func testRegistry(t *testing.T, imageName string) {
	var r inspector.Registry
	inspect, err := r.Inspect(context.Background(), imageName)
	assert.NoError(t, err)
	assert.NotNil(t, inspect)
	assert.NotNil(t, inspect.Config)
	assert.Equal(t, "aws", inspect.Config.Labels["io.whalebrew.name"])
}

func TestRegistry(t *testing.T) {
	testRegistry(t, "whalebrew/awscli")
	testRegistry(t, "whalebrew/awscli:latest")
	testRegistry(t, "docker.io/whalebrew/awscli:latest")
	testRegistry(t, "docker.io/whalebrew/awscli@sha256:bf82c73991d98b1e5f6999928176a94e7597ccfb2a366e49aad5122d3a1da9f8")
}
