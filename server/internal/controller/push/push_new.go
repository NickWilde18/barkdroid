package push

import (
	"barkdroid/internal/provider"
	"barkdroid/internal/store"
)

// Controller 推送控制器。
type Controller struct {
	store     store.Store
	providers map[string]provider.Provider
}

// New 创建推送控制器。
func New(st store.Store, providers map[string]provider.Provider) *Controller {
	return &Controller{
		store:     st,
		providers: providers,
	}
}
