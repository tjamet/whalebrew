package inspector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	dist "github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/distribution"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/registry"
	"github.com/whalebrew/whalebrew/client"
)

// Registry allows inspecting an image directly from a registry without downloading it
type Registry struct{}

// GetRepository returns a repository from the registry
// from https://github.com/moby/moby/blob/ad1b781e44fa1e44b9e654e5078929aec56aed66/daemon/images/image_pull.go#L91
// using the code above leads to dependency hell:
// github.com/whalebrew/whalebrew/* imports
//         github.com/docker/docker/daemon/images imports
//         github.com/docker/docker/container imports
//         github.com/docker/swarmkit/agent/exec imports
//         github.com/docker/swarmkit/api imports
//         google.golang.org/grpc/transport: module google.golang.org/grpc@latest (v1.24.0) found, but does not contain package google.golang.org/grpc/transport
func GetRepository(ctx context.Context, registryService registry.Service, ref reference.Named, authConfig *types.AuthConfig) (dist.Repository, bool, error) {
	// get repository info
	repoInfo, err := registryService.ResolveRepository(ref)
	if err != nil {
		return nil, false, errdefs.InvalidParameter(err)
	}
	// makes sure name is not empty or `scratch`
	if err := distribution.ValidateRepoName(repoInfo.Name); err != nil {
		return nil, false, errdefs.InvalidParameter(err)
	}

	// get endpoints
	endpoints, err := registryService.LookupPullEndpoints(reference.Domain(repoInfo.Name))
	if err != nil {
		return nil, false, err
	}

	// retrieve repository
	var (
		confirmedV2 bool
		repository  dist.Repository
		lastError   error
	)

	for _, endpoint := range endpoints {
		if endpoint.Version == registry.APIVersion1 {
			continue
		}

		repository, confirmedV2, lastError = distribution.NewV2Repository(ctx, repoInfo, endpoint, nil, authConfig, "pull")
		if lastError == nil && confirmedV2 {
			break
		}
	}
	return repository, confirmedV2, lastError
}

func (r *Registry) Inspect(ctx context.Context, imageName string) (*types.ImageInspect, error) {
	registryService, err := registry.NewService(registry.ServiceOptions{})
	if err != nil {
		return nil, err
	}
	anyRef, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return nil, err
	}
	ref, ok := anyRef.(reference.Named)
	if !ok {
		return nil, fmt.Errorf("unsupported reference type: %T for image %s", ref, imageName)
	}
	ref = reference.TagNameOnly(ref)
	repo, _, err := GetRepository(ctx, registryService, ref, client.AuthConfig(ctx, ref))
	if err != nil {
		return nil, err
	}
	manSvc, err := repo.Manifests(ctx)
	if err != nil {
		return nil, err
	}
	var manifest dist.Manifest
	//https://github.com/moby/moby/blob/ad1b781e44fa1e44b9e654e5078929aec56aed66/distribution/pull_v2.go#L331
	if digested, isDigested := ref.(reference.Canonical); isDigested {
		manifest, err = manSvc.Get(ctx, digested.Digest())
	} else if tagged, isTagged := ref.(reference.NamedTagged); isTagged {
		manifest, err = manSvc.Get(ctx, "", dist.WithTag(tagged.Tag()))
	}
	if manifest == nil {
		return nil, fmt.Errorf("unable to get manifest for image %s", imageName)
	}
	if m, ok := manifest.(*schema2.DeserializedManifest); ok {
		configJSON, err := repo.Blobs(ctx).Get(ctx, m.Target().Digest)
		if err != nil {
			return nil, err
		}
		r := types.ImageInspect{}
		return &r, json.NewDecoder(bytes.NewReader(configJSON)).Decode(&r)
	}
	return nil, fmt.Errorf("unsupported manifest type %T", manifest)

}
