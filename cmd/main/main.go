package main

import (
	"context"
	"log"
	"net"
	"os"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/soulmate-dating/media/internal/adapters/s3"
	"github.com/soulmate-dating/media/internal/app"
	"github.com/soulmate-dating/media/internal/graceful"
	grpcSvc "github.com/soulmate-dating/media/internal/ports/grpc"
)

const (
	grpcPort = ":8082"
)

func main() {
	ctx := context.Background()

	s3Client, err := s3.NewClient("minio:9000", "accesskey", "secretkey", "media", false)
	if err != nil {
		log.Fatal(err.Error())
	}
	appSvc := app.NewApp(s3Client)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterMediaServiceServer(grpcServer, svc)

	eg, ctx := errgroup.WithContext(ctx)

	sigQuit := make(chan os.Signal, 1)
	eg.Go(graceful.CaptureSignal(ctx, sigQuit))
	// run grpc server
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, lis, grpcServer))

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}
	log.Println("servers were successfully shutdown")
}
