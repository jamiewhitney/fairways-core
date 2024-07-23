package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"github.com/jamiewhitney/fairways-core/pkg/transport"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	catalog_pb "github.com/jamiewhitney/fairways-core/protobufs/catalog"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"github.com/jamiewhitney/fairways-core/protobufs/tee_time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"os"
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

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:     true,
			EmitDefaultValues: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pricing_pb.RegisterPricingServiceHandlerFromEndpoint(ctx, mux, os.Getenv("PRICING_ADDR"), opts)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = catalog_pb.RegisterCatalogServiceHandlerFromEndpoint(ctx, mux, os.Getenv("CATALOG_ADDR"), opts)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = tee_time.RegisterTeeTimeServiceHandlerFromEndpoint(ctx, mux, os.Getenv("TEE_TIME_ADDR"), opts)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	err = booking_pb.RegisterBookingServiceHandlerFromEndpoint(ctx, mux, os.Getenv("BOOKINGS_ADDR"), opts)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	srv, err := transport.New("8081")
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = mux.HandlePath(http.MethodGet, "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"status\":\"ok\"}"))
	})
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		if err := srv.ServeHTTP(ctx, &http.Server{Handler: mux}); err != nil {
			logger.Fatalf(err.Error())
		}
	}()

	logger.Infof("http server started on port :%s", "8081")

	<-ctx.Done()

	logger.Infof("shutting down server")
}
