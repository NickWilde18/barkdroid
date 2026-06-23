package main

import (
	"log"

	"barkdroid/internal/handler"
	"barkdroid/internal/provider"
	"barkdroid/internal/store"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()

	// --- config ---
	dbPath := g.Cfg().MustGet(ctx, "db.path", "data/barkdroid.db").String()
	serverAddr := g.Cfg().MustGet(ctx, "server.address", ":8080").String()

	// --- store ---
	st, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}
	defer st.Close()
	log.Printf("[INFO] SQLite store ready at %s", dbPath)

	// --- providers ---
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

	// --- routes ---
	h := handler.New(st, providers)

	s := g.Server()

	// Bark-compatible endpoints
	s.BindHandler("/:key/:title/:body", h.PushTitleBody)
	s.BindHandler("/:key/:body", h.PushBody)

	// Additional endpoints
	s.BindHandler("POST:/push", h.PushPost)
	s.BindHandler("POST:/register", h.RegisterDevice)

	// Health check
	s.BindHandler("GET:/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"status": "ok"})
	})

	s.SetAddr(serverAddr)

	log.Printf("[INFO] barkdroid server starting on %s", serverAddr)
	s.Run()
}
