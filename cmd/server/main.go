package main

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"OpenCodeGoPool/internal/api"
	"OpenCodeGoPool/internal/config"
	"OpenCodeGoPool/internal/cpa"
	"OpenCodeGoPool/internal/frontend"
	"OpenCodeGoPool/internal/model"
	"OpenCodeGoPool/internal/scheduler"
	"OpenCodeGoPool/internal/scraper"
	"OpenCodeGoPool/internal/store"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	db, err := store.NewSQLite(cfg.Database.Path)
	if err != nil {
		slog.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	sc := scraper.NewClient(cfg.Scraper.Timeout)
	cpaSyncer := cpa.NewSyncer(db)

	cpaSyncer.InitDefaults(context.Background(), model.CPASettings{
		Endpoint:     cfg.CPA.Endpoint,
		BearerToken:  cfg.CPA.BearerToken,
		ProviderName: cfg.CPA.ProviderName,
		BaseURL:      cfg.CPA.BaseURL,
		Models:       model.DefaultCPAModels,
	})

	sched := scheduler.New(db, sc, cpaSyncer, cfg.Scraper.Interval)
	sched.Start()
	defer sched.Stop()

	var frontendFS fs.FS
	distFS, err := fs.Sub(frontend.DistFS, "dist")
	if err == nil {
		frontendFS = distFS
	}

	srv := api.NewServer(db, sched, cpaSyncer, cfg.Server.Password)
	httpServer := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: srv.Router(frontendFS),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("server starting", "address", cfg.Server.Address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(shutdownCtx)
}
