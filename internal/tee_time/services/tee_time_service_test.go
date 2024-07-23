package services

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	bookings "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	mocks "github.com/jamiewhitney/fairways-core/testing/mocks/repositories"
	"github.com/jamiewhitney/fairways-core/testing/mocks/services"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

var (
	lis    *bufconn.Listener
	logger = logging.NewLoggerFromEnv()
)

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	lis = bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pricing.RegisterPricingServiceServer(server, &services.PricingService{
		Price: 6.0,
	})
	bookings.RegisterBookingServiceServer(server, &services.BookingService{})
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	logger.Out = io.Discard
}

func TestTeeTimeService_(t *testing.T) {
	tests := []struct {
		name          string
		courseId      int64
		date          time.Time
		golfers       int64
		mockPrice     float64
		mockError     error
		expectedError bool
	}{
		{
			name:          "valid request",
			courseId:      10046,
			date:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			mockPrice:     4.0,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "schedule/course does not exist",
			courseId:      10047,
			date:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			mockPrice:     0,
			mockError:     errors.New("test error"),
			expectedError: true,
		},
	}

	ctx := context.Background()
	ctx = logging.WithLogger(ctx, logger)

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	defer conn.Close()

	redisClient, _ := redismock.NewClientMock()
	pricingServiceClient := pricing.NewPricingServiceClient(conn)
	bookingServiceClient := bookings.NewBookingServiceClient(conn)

	ttc := &TeeTimeService{
		teeTimeRepository: &mocks.MockTeeTimeRepository{},
		cache:             redisClient,
		pricingService:    pricingServiceClient,
		bookingService:    bookingServiceClient,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ttc.GetTeeTimes(ctx, tt.courseId, tt.date.Format(time.RFC3339), tt.golfers)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
