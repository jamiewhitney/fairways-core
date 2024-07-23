package main

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/pricing"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/setup"
	"github.com/jamiewhitney/fairways-core/pkg/transport"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logger := logging.NewLoggerFromEnv()

	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatal("application panic", "panic", r)
		}
	}()

	var config pricing.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	ps := pricing.New(env, &config)

	grpcServer := transport.NewGRPCServer()

	pricing_pb.RegisterPricingServiceServer(grpcServer, &ps)

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
