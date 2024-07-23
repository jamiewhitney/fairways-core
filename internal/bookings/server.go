package bookings

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/bookings/repository"
	"github.com/jamiewhitney/fairways-core/internal/bookings/services"
	"github.com/jamiewhitney/fairways-core/pkg/environment"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	"time"
)

type Server struct {
	booking_pb.UnimplementedBookingServiceServer
	bookingService *services.BookingService
}

func New(env *environment.Environment, config *Config) Server {
	return Server{
		bookingService: services.NewBookingService(repository.New(env.Database().Pool)),
	}
}

func (s *Server) GetConfirmedBookings(ctx context.Context, in *booking_pb.GetConfirmedBookingsRequest) (*booking_pb.GetConfirmedBookingResponse, error) {
	parsedTime, err := time.Parse(time.RFC3339, in.Datetime)
	if err != nil {
		return nil, err
	}

	return s.bookingService.GetConfirmedBookings(ctx, in.CourseId, parsedTime)

}

func (s *Server) GetBookings(ctx context.Context, in *booking_pb.GetBookingsRequest) (*booking_pb.GetBookingsResponse, error) {
	result, err := s.bookingService.GetBookings(ctx, in.UserId)
	if err != nil {
		return nil, err
	}

	return result, nil
}
