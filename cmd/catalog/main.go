package main

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/catalog"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/setup"
	"github.com/jamiewhitney/fairways-core/pkg/transport"
	catalog_pb "github.com/jamiewhitney/fairways-core/protobufs/catalog"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logger := logging.NewLoggerFromEnv()

	ctx = logging.WithLogger(ctx, logger)

	var config catalog.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	ps := catalog.New(env, &config)
	grpcServer := transport.NewGRPCServer()

	catalog_pb.RegisterCatalogServiceServer(grpcServer, &ps)

	tt, err := transport.New(config.Port)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	go func() {
		if err := tt.ServeGRPC(ctx, grpcServer); err != nil {
			logger.Fatalf(err.Error())
		}
	}()

	logger.Infof("grpc server started on port :%s", config.Port)

	<-ctx.Done()

	logger.Info("shutting down server")
}
