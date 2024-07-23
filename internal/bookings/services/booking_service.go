package services

import (
	"context"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/bookings/repository"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	"time"
)

type BookingService struct {
	bookingRepository repository.Querier
}

func NewBookingService(bookingRepository repository.Querier) *BookingService {
	return &BookingService{
		bookingRepository: bookingRepository,
	}
}

func (bs *BookingService) GetBookings(ctx context.Context, userId string) (*booking_pb.GetBookingsResponse, error) {
	result, err := bs.bookingRepository.ListUserBookings(ctx, userId)
	if err != nil {
		return nil, err
	}

	var bookings []*booking_pb.Booking
	for _, booking := range result {
		bookings = append(bookings, &booking_pb.Booking{
			Id:              booking.ID,
			UserId:          booking.UserID,
			CourseId:        booking.CourseID,
			Golfers:         booking.Golfers,
			Datetime:        booking.Datetime.String(),
			BookingId:       booking.BookingID,
			StripePaymentId: booking.StripePaymentID,
		})
	}
	return &booking_pb.GetBookingsResponse{
		Bookings: bookings,
	}, nil
}

func (bs *BookingService) GetConfirmedBookings(ctx context.Context, courseId int64, datetime time.Time) (*booking_pb.GetConfirmedBookingResponse, error) {
	result, err := bs.bookingRepository.GetConfirmedBookingsByDateAndCourse(ctx, &repository.GetConfirmedBookingsByDateAndCourseParams{
		CourseID: courseId,
		Datetime: datetime,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	var bookings []*booking_pb.Booking
	for _, booking := range result {
		bookings = append(bookings, &booking_pb.Booking{
			Id:        booking.ID,
			Datetime:  booking.Datetime.String(),
			Golfers:   booking.Golfers,
			BookingId: booking.BookingID,
			CourseId:  booking.CourseID,
		})
	}
	return &booking_pb.GetConfirmedBookingResponse{
		Bookings: bookings,
	}, nil
}
