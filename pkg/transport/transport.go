package transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"strconv"
	"time"
)

var defaultServerOpts = []grpc.ServerOption{grpc.UnaryInterceptor(unaryLoggingInterceptor), grpc.StreamInterceptor(streamLoggingInterceptor)}
var defaultClientOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

type Server struct {
	ip       string
	port     string
	listener net.Listener
}

func New(port string) (*Server, error) {
	if port == "" {
		return nil, errors.New("port not set")
	}

	addr := fmt.Sprintf(":" + port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func NewGRPCServer(opts ...[]grpc.ServerOption) *grpc.Server {
	var serverOpts []grpc.ServerOption

	for _, opt := range defaultServerOpts {
		serverOpts = append(serverOpts, opt)
	}

	for _, opt := range opts {
		serverOpts = append(serverOpts, opt...)
	}

	srv := grpc.NewServer(serverOpts...)
	hs := health.NewServer()
	healthpb.RegisterHealthServer(srv, hs)

	return srv
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	logger := logging.FromContext(ctx)
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		logger.Info("server.Serve: context closed")
		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		logger.Info("server.Serve: shutting down")
		errCh <- srv.Shutdown(shutdownCtx)
	}()

	srv.Handler = NewRequestID(LoggingMiddleware(srv.Handler))

	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Info("server.Serve: serving stopped")

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) ServeGRPC(ctx context.Context, srv *grpc.Server) error {
	logger := logging.FromContext(ctx)
	go func() {
		<-ctx.Done()
		logger.Debugf("server.Serve: context closed")
		logger.Debugf("server.Serve: shutting down")
		srv.GracefulStop()
	}()

	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}

func NewClient(ctx context.Context, addr string, opts []grpc.DialOption) (*grpc.ClientConn, error) {
	var clientOpts []grpc.DialOption

	for _, opt := range defaultClientOpts {
		clientOpts = append(clientOpts, opt)
	}

	for _, opt := range opts {
		clientOpts = append(clientOpts, opt)
	}

	logger := logging.FromContext(ctx)
	conn, err := grpc.NewClient(addr, clientOpts...)
	if err != nil {
		logger.Errorf("failed to dial: %v", err)
		return nil, err
	}
	return conn, nil
}
