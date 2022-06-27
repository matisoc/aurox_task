package app

import (
	"aurox_task/internal/modules"
	"context"
)

func Run(ctx context.Context, settings map[string]interface{}) {
	modules.NewSiteMapGenerator(ctx, settings)
}
