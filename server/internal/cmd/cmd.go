package cmd

import (
	"context"
	"fmt"
	"log"

	pushctrl "barkdroid/internal/controller/push"
	"barkdroid/internal/provider"
	"barkdroid/internal/store"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

// Main 是 barkdroid 服务端入口命令。
var Main = gcmd.Command{
	Name:  "barkdroid",
	Usage: "barkdroid",
	Brief: "Bark-compatible Android push server",
	Func: func(ctx context.Context, parser *gcmd.Parser) error {
		return run(ctx)
	},
}

func run(ctx context.Context) error {
	dbPath := g.Cfg().MustGet(ctx, "db.path", "data/barkdroid.db").String()
	serverAddr := g.Cfg().MustGet(ctx, "server.address", ":8080").String()

	st, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		return fmt.Errorf("init store: %w", err)
	}
	defer st.Close()
	log.Printf("[INFO] SQLite store ready at %s", dbPath)

	providers := make(map[string]provider.Provider)

	if g.Cfg().MustGet(ctx, "jpush.enabled", false).Bool() {
		providers["jpush"] = provider.NewJPush(provider.JPushConfig{
			AppKey:       g.Cfg().MustGet(ctx, "jpush.app_key").String(),
			MasterSecret: g.Cfg().MustGet(ctx, "jpush.master_secret").String(),
		})
		log.Println("[INFO] JPush provider enabled")
	}

	if g.Cfg().MustGet(ctx, "bark.enabled", false).Bool() {
		providers["bark"] = provider.NewBark(provider.BarkConfig{
			BaseURL: g.Cfg().MustGet(ctx, "bark.base_url", "https://api.day.app").String(),
			Key:     g.Cfg().MustGet(ctx, "bark.key").String(),
		})
		log.Println("[INFO] Bark forward provider enabled")
	}

	ctrl := pushctrl.New(st, providers)

	s := g.Server()

	s.BindHandler("/:key/:title/:body", ctrl.PushTitleBody)
	s.BindHandler("/:key/:body", ctrl.PushBody)
	s.BindHandler("POST:/push", ctrl.PushPost)
	s.BindHandler("POST:/register", ctrl.RegisterDevice)
	s.BindHandler("GET:/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"status": "ok"})
	})

	s.SetAddr(serverAddr)

	log.Printf("[INFO] barkdroid server starting on %s", serverAddr)
	s.Run()
	return nil
}
