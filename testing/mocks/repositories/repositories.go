package mocks

import (
	"context"
	"fmt"
	"github.com/jamiewhitney/fairways-core/internal/tee_time/repository"
	"time"
)

type MockTeeTimeRepository struct {
}

func (t MockTeeTimeRepository) InsertSchedule(ctx context.Context, arg *repository.InsertScheduleParams) error {
	//TODO implement me
	panic("implement me")
}

func (t MockTeeTimeRepository) GetSchedules(ctx context.Context, courseID int64) ([]repository.GetSchedulesRow, error) {
	//TODO implement me
	panic("implement me")
}

func (t MockTeeTimeRepository) GetOverrides(ctx context.Context, arg *repository.GetOverridesParams) ([]repository.GetOverridesRow, error) {
	if arg.CourseID == 10046 {

		return []repository.GetOverridesRow{
			{
				StartTime: time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC),
				EndTime:   time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC),
				Blocked:   true,
			},
		}, nil
	} else {
		return nil, fmt.Errorf("course not found")
	}
}

func (t MockTeeTimeRepository) GetSchedule(ctx context.Context, arg *repository.GetScheduleParams) (repository.GetScheduleRow, error) {
	return repository.GetScheduleRow{
		CourseID:   10046,
		StartTime:  time.Date(0, 0, 0, 7, 0, 0, 0, time.UTC),
		EndTime:    time.Date(0, 0, 0, 19, 0, 0, 0, time.UTC),
		Occurrence: 12,
		Day:        1,
		Buffer:     12,
	}, nil
}
