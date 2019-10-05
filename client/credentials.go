package client

import (
	"context"
	"os"

	dockerConfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/trust"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/distribution/reference"
)

func resolveAuthConfig(ctx context.Context, index *registrytypes.IndexInfo) types.AuthConfig {
	configFile := dockerConfig.LoadDefaultConfigFile(os.Stderr)
	a, _ := configFile.GetAuthConfig(index.Name)
	return types.AuthConfig(a)
}

// AuthConfig retrieves the authentication configuration for the image
func AuthConfig(ctx context.Context, ref reference.Reference) *types.AuthConfig {
	refAndAuth, err := trust.GetImageReferencesAndAuth(ctx, nil, resolveAuthConfig, ref.String())
	if err != nil {
		return nil
	}
	return refAndAuth.AuthConfig()
}
