package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/tee_time/repository"
	databasesql "github.com/jamiewhitney/fairways-core/pkg/database/mysql"
	"github.com/jamiewhitney/fairways-core/pkg/logging"
	booking_pb "github.com/jamiewhitney/fairways-core/protobufs/booking"
	pricing_pb "github.com/jamiewhitney/fairways-core/protobufs/pricing"
	"github.com/jamiewhitney/fairways-core/protobufs/tee_time"
	"github.com/redis/go-redis/v9"
	"io"
	"os"
	"sync"
	"time"
)

type TeeTimeService struct {
	teeTimeRepository repository.Querier
	cache             redis.Cmdable

	pricingService pricing_pb.PricingServiceClient
	bookingService booking_pb.BookingServiceClient
}

func NewTeeTimeService(db *databasesql.DB, cache redis.Cmdable, pricingService pricing_pb.PricingServiceClient, bookingService booking_pb.BookingServiceClient) *TeeTimeService {
	return &TeeTimeService{
		teeTimeRepository: repository.New(db.Pool),
		cache:             cache,
		pricingService:    pricingService,
		bookingService:    bookingService,
	}
}

func (tts *TeeTimeService) GetTeeTimes(ctx context.Context, courseId int64, date string, golfers int64) (*tee_time.GetTeeTimesResponse, error) {
	logger := logging.FromContext(ctx)

	parsedTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, err
	}

	key := cacheKey(courseId, parsedTime, golfers)

	// Try to get cached data
	if os.Getenv("DISABLE_CACHE") != "false" {
		cachedData, err := tts.cache.Get(ctx, key).Result()
		if err == nil {
			var teetimes []*tee_time.TeeTime
			if jsonErr := json.Unmarshal([]byte(cachedData), &teetimes); jsonErr == nil {
				return &tee_time.GetTeeTimesResponse{
					Teetimes: teetimes,
				}, nil
			}
		}
	}
	// get schedule for the course
	schedules, err := tts.getSchedule(ctx, courseId, parsedTime)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// Generate times list from schedules
	timesList := generateTimesList(ctx, schedules, parsedTime)

	// get overrides for the course
	overrides, err := tts.getOverrides(ctx, courseId, parsedTime)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// get confirmed bookings for the course
	confirmedBookings, err := tts.getBookings(ctx, courseId, parsedTime)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	teetimes, err := tts.getPricingAndGenerateTeeTimes(ctx, courseId, golfers, timesList, overrides, confirmedBookings, time.Duration(schedules.Buffer)*time.Minute)
	if err != nil {
		return nil, err
	}

	//Cache the result
	if os.Getenv("DISABLE_CACHE") != "true" {
		data, err := json.Marshal(teetimes)
		if err == nil {
			tts.cache.Set(ctx, key, data, 10*time.Minute)
		} else {
			logger.Warn(err)
		}
	}

	return &tee_time.GetTeeTimesResponse{
		Teetimes: teetimes,
	}, nil
}

// generateTimesList generates a list of times based on the schedules
func generateTimesList(ctx context.Context, schedule *repository.GetScheduleRow, date time.Time) []time.Time {
	logger := logging.FromContext(ctx)

	logger.Debug("generating times from schedule")

	var timesList []time.Time
	interval := time.Duration(schedule.Occurrence) * time.Minute

	timesList = append(timesList, generateTimes(schedule.StartTime, schedule.EndTime, interval, date)...)

	logger.Debugf("times generated %s", timesList)
	return timesList
}

func (tts *TeeTimeService) PreloadTeeTimes(ctx context.Context, courseDates map[int64][]time.Time) error {
	logger := logging.FromContext(ctx)
	for courseId, dates := range courseDates {
		for _, date := range dates {
			key := cacheKey(courseId, date, 1)

			// Check if data already exists in cache
			exists, err := tts.cache.Exists(ctx, key).Result()
			if err != nil {
				logger.Errorf("Error checking cache for course %d on date %s: %v", courseId, date.Format("2006-01-02"), err)
				continue
			}

			if exists > 0 {
				logger.Infof("Cache already exists for course %d on date %s", courseId, date.Format("2006-01-02"))
				continue
			}

			// If not cached, get tee times and cache them
			teetimes, err := tts.GetTeeTimes(ctx, courseId, date.String(), 1)
			if err != nil {
				logger.Errorf("Error preloading tee times for course %d on date %s: %v", courseId, date.Format("2006-01-02"), err)
				continue
			}

			logger.Infof("Preloaded %d tee times for course %d on date %s", len(teetimes.Teetimes), courseId, date.Format("2006-01-02"))
		}
	}
	return nil
}

func isBlocked(ctx context.Context, timeToCheck time.Time, overrides *[]repository.GetOverridesRow, bookings []*booking_pb.Booking, buffer time.Duration) bool {
	logger := logging.FromContext(ctx)

	for _, override := range *overrides {
		if (timeToCheck.After(override.StartTime) && timeToCheck.Before(override.EndTime)) ||
			timeToCheck.Equal(override.StartTime) ||
			timeToCheck.Equal(override.EndTime) {
			return true
		}
	}

	for _, booking := range bookings {
		bookingTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", booking.Datetime)
		if err != nil {
			logger.Error(err)
			continue
		}
		if timeToCheck.Equal(bookingTime) ||
			(timeToCheck.After(bookingTime.Add(-buffer)) && timeToCheck.Before(bookingTime.Add(buffer))) {
			return true
		}
	}

	return false
}

func (tts *TeeTimeService) InvalidateCache(ctx context.Context, courseId int64, date string, basePrice bool) error {
	var key string
	if basePrice {
		key = fmt.Sprintf("*%v", courseId)
	}

	key = fmt.Sprintf("*%v:%s*", courseId, date)

	keys, err := tts.cache.Keys(ctx, key).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := tts.cache.Del(ctx, key); err != nil {
			return err.Err()
		}
	}

	return nil
}

func cacheKey(courseId int64, date time.Time, golfers int64) string {
	return fmt.Sprintf("teetimes:%d:%s:%v", courseId, date.Format("2006-01-02"), golfers)
}

func generateTimes(startTime, endTime time.Time, interval time.Duration, date time.Time) []time.Time {
	var times []time.Time
	current := startTime

	for current.Before(endTime) || current.Equal(endTime) {
		times = append(times, time.Date(date.Year(), date.Month(), date.Day(), current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), current.Location()))
		current = current.Add(interval)
	}

	return times
}

func (tts *TeeTimeService) getSchedule(ctx context.Context, courseId int64, date time.Time) (*repository.GetScheduleRow, error) {
	logger := logging.FromContext(ctx)
	logger.Debugf("retrieving schedule course_id: %v date: %v", courseId, date)

	schedules, err := tts.teeTimeRepository.GetSchedule(ctx, &repository.GetScheduleParams{
		CourseID: courseId,
		Day:      int64(date.Weekday()),
	})
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}
	logger.Debugf("retrieved schedule %+v", schedules)

	return &schedules, err
}

func (tts *TeeTimeService) getOverrides(ctx context.Context, courseId int64, date time.Time) (*[]repository.GetOverridesRow, error) {
	logger := logging.FromContext(ctx)

	logger.Debugf("retrieving overrides course_id: %v date: %v", courseId, date)

	overrides, err := tts.teeTimeRepository.GetOverrides(ctx, &repository.GetOverridesParams{
		CourseID: courseId,
		Date:     date.Truncate(24 * time.Hour),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule overrides: %w", err)
	}

	logger.Debugf("retrieved overrides %+v", overrides)

	return &overrides, nil
}

func (tts *TeeTimeService) getBookings(ctx context.Context, courseId int64, date time.Time) ([]*booking_pb.Booking, error) {
	logger := logging.FromContext(ctx)

	logger.Debugf("retrieving confirmed bookings course_id: %v date: %v", courseId, date)

	confirmedBookings, err := tts.bookingService.GetConfirmedBookings(ctx, &booking_pb.GetConfirmedBookingsRequest{
		CourseId: courseId,
		Datetime: date.Format(time.RFC3339),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmed bookings: %w", err)
	}
	logger.Debugf("retrieved confirmed bookings %+v", confirmedBookings)

	return confirmedBookings.Bookings, nil
}

func (tts *TeeTimeService) getPricingAndGenerateTeeTimes(ctx context.Context, courseId int64, golfers int64, timesList []time.Time, overrides *[]repository.GetOverridesRow, confirmedBookings []*booking_pb.Booking, buffer time.Duration) ([]*tee_time.TeeTime, error) {
	stream, err := tts.pricingService.GetPriceStream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get price stream: %w", err)
	}

	var teetimes []*tee_time.TeeTime
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- err
				return
			}

			mu.Lock()
			teetimes = append(teetimes, &tee_time.TeeTime{
				CourseId:      uint64(courseId),
				Datetime:      resp.Datetime,
				Price:         resp.Price,
				OriginalPrice: resp.OriginalPrice,
				Discounted:    resp.Discounted,
				Available:     true,
			})
			mu.Unlock()
		}
		close(errChan)
	}()

	for _, t := range timesList {
		if !isBlocked(ctx, t, overrides, confirmedBookings, buffer) {
			if err := stream.Send(&pricing_pb.GetPriceRequest{
				CourseId: courseId,
				Datetime: t.Format(time.RFC3339),
				Golfers:  golfers,
			}); err != nil {
				return nil, fmt.Errorf("failed to send price request: %w", err)
			}
		}
	}

	if err := stream.CloseSend(); err != nil {
		return nil, fmt.Errorf("failed to closa	e send stream: %w", err)
	}

	wg.Wait()

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return teetimes, nil
}

func (tts *TeeTimeService) GetSchedule(ctx context.Context, courseId int64) (*tee_time.GeeTeeTimeScheduleResponse, error) {
	logger := logging.FromContext(ctx)

	schedules, err := tts.teeTimeRepository.GetSchedules(ctx, courseId)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if schedules == nil {
		return nil, fmt.Errorf("no schedules found for course %d", courseId)
	}

	var teeTimeSchedules []*tee_time.TeeTimeSchedule
	for _, schedule := range schedules {
		teeTimeSchedules = append(teeTimeSchedules, &tee_time.TeeTimeSchedule{
			CourseId:   schedule.CourseID,
			Day:        schedule.Day,
			StartTime:  schedule.StartTime.Format(time.TimeOnly),
			EndTime:    schedule.EndTime.Format(time.TimeOnly),
			Occurrence: schedule.Occurrence,
		})
	}

	return &tee_time.GeeTeeTimeScheduleResponse{
		Schedule: teeTimeSchedules,
	}, nil

}

func (tts *TeeTimeService) CreateSchedule(ctx context.Context, courseId int64, day int64, startTime string, endTime string, occurrence int64) error {
	logger := logging.FromContext(ctx)

	st, err := time.Parse(time.TimeOnly, startTime)
	if err != nil {
		logger.Error(err)
		return err
	}

	et, err := time.Parse(time.TimeOnly, endTime)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = tts.teeTimeRepository.InsertSchedule(ctx, &repository.InsertScheduleParams{
		CourseID:   courseId,
		Day:        day,
		StartTime:  st.AddDate(2020, 0, 0),
		EndTime:    et.AddDate(2020, 0, 0),
		Occurrence: occurrence,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
