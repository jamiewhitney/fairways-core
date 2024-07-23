package services

import (
	"context"
	bookings "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"io"
	"strconv"
	"time"
)

type MockClientStream struct {
	grpc.ClientStream
	Requests  []*pricing.GetPriceRequest
	Responses []*pricing.GetPriceResponse
	RecvIndex int
	mock.Mock
}

func (m *MockClientStream) Send(req *pricing.GetPriceRequest) error {
	m.Requests = append(m.Requests, req)
	return nil
}

func (m *MockClientStream) Recv() (*pricing.GetPriceResponse, error) {
	if m.RecvIndex < len(m.Responses) {
		resp := m.Responses[m.RecvIndex]
		m.RecvIndex++
		return resp, nil
	}
	return nil, io.EOF
}

type PricingService struct {
	pricing.UnimplementedPricingServiceServer
	Price float64
}

type BookingService struct {
	bookings.UnimplementedBookingServiceServer
	Price float64
}

func (s *BookingService) GetConfirmedBookings(ctx context.Context, in *bookings.GetConfirmedBookingsRequest) (*bookings.GetConfirmedBookingResponse, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", in.Datetime)
	if err != nil {
		return nil, err
	}

	return &bookings.GetConfirmedBookingResponse{
		Bookings: []*bookings.Booking{
			{
				CourseId: in.CourseId,
				Datetime: time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 9, 11, 0, 0, time.UTC).String(),
			},
			{
				CourseId: in.CourseId,
				Datetime: time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 9, 24, 0, 0, time.UTC).String(),
			},
		},
	}, nil
}

func (s *PricingService) GetPriceStream(stream pricing.PricingService_GetPriceStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		response := &pricing.GetPriceResponse{
			Price:    s.Price,
			CourseId: strconv.FormatInt(req.CourseId, 10),
			Datetime: req.Datetime,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func (s *PricingService) GetPrice(ctx context.Context, in *pricing.GetPriceRequest) (*pricing.GetPriceResponse, error) {
	return &pricing.GetPriceResponse{
		Price:    s.Price,
		CourseId: strconv.FormatInt(in.CourseId, 10),
		Datetime: in.Datetime,
	}, nil
}
