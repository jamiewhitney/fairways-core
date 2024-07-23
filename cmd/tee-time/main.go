package main

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/tee_time"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/setup"
	"github.com/jamiewhitney/fairways-core/pkg/transport"
	tee_time_pb "github.com/jamiewhitney/fairways-core/protobufs/tee_time"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer done()

	logger := logging.NewLoggerFromEnv()

	ctx = logging.WithLogger(ctx, logger)

	var config tee_time.Config
	env, err := setup.Setup(ctx, &config)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	config.Port = "3008"

	tts := tee_time.New(env, &config)

	grpcServer := transport.NewGRPCServer()

	tee_time_pb.RegisterTeeTimeServiceServer(grpcServer, &tts)

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
