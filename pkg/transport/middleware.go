package transport

import (
	"context"
	"github.com/google/uuid"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	"google.golang.org/grpc"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.FromContext(r.Context())
		logger.Println(r.Method, r.RequestURI, r.Proto, r.Host, r.RemoteAddr, r.Header.Get("X-Request-Id"))
		next.ServeHTTP(w, r)
	})
}

func NewRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New().String()
		r.Header.Set("X-Request-Id", uuid)
		w.Header().Set("X-Request-Id", uuid)
		next.ServeHTTP(w, r)
	})
}

func unaryLoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger := logging.FromContext(ctx)
	logger.Infof(info.FullMethod)
	return handler(ctx, req)
}

func streamLoggingInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger := logging.FromContext(ss.Context())
	logger.Info(info.FullMethod)
	return handler(srv, ss)
}
