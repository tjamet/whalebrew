package inspector

import (
	"context"

	"github.com/docker/docker/api/types"
)

type Inspector interface {
	Inspect (context.Context, string)  (*types.ImageInspect, error)
}