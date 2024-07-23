package tee_time

import (
	"context"
	"github.com/jamiewhitney/fairways-core/internal/tee_time/services"
	"github.com/jamiewhitney/fairways-core/pkg/environment"
	"github.com/jamiewhitney/fairways-core/protobufs/tee_time"
)

type Server struct {
	tee_time.UnimplementedTeeTimeServiceServer

	teeTimeService *services.TeeTimeService
}

func New(env *environment.Environment, config *Config) Server {
	return Server{
		teeTimeService: services.NewTeeTimeService(env.Database(), env.Cache(), env.PricingClient(), env.BookingClient()),
	}
}

func (s *Server) GetTeeTimeByDateTime(context.Context, *tee_time.GetTeeTimeDateRequest) (*tee_time.GetTeeTimeResponse, error) {
	panic("implement me")
}
func (s *Server) GetTeeTimes(ctx context.Context, in *tee_time.GetTeeTimesRequest) (*tee_time.GetTeeTimesResponse, error) {
	return s.teeTimeService.GetTeeTimes(ctx, in.CourseId, in.Date, in.Golfers)
}
func (s *Server) UpdateTeeTimeAvailability(context.Context, *tee_time.TeeTime) (*tee_time.TeeTime, error) {
	panic("implement me")
}
func (s *Server) InvalidateCache(context.Context, *tee_time.InvalidateCacheRequest) (*tee_time.InvalidateCacheResponse, error) {
	panic("implement me")
}

func (s *Server) GetTeeTimeSchedules(ctx context.Context, in *tee_time.GetTeeTimeScheduleRequest) (*tee_time.GeeTeeTimeScheduleResponse, error) {
	return s.teeTimeService.GetSchedule(ctx, in.CourseId)
}
func (s *Server) CreateTeeTimeSchedule(ctx context.Context, in *tee_time.CreateTeeTimeScheduleRequest) (*tee_time.CreateTeeTimeScheduleResponse, error) {
	return &tee_time.CreateTeeTimeScheduleResponse{Created: true}, s.teeTimeService.CreateSchedule(ctx, in.CourseId, in.Day, in.StartTime, in.EndTime, in.Occurrence)
}
