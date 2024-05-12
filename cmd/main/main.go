package main

import (
	"context"
	"github.com/soulmate-dating/media/internal/ports/grpc"
	"log"

	"github.com/soulmate-dating/media/internal/app"
	"github.com/soulmate-dating/media/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	appSvc := app.New(ctx, cfg)
	grpc.Run(ctx, cfg, appSvc)
}
