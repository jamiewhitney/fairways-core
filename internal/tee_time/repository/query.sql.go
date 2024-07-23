// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package repository

import (
	"context"
	"time"
)

const getOverrides = `-- name: GetOverrides :many
SELECT start_time, end_time, blocked
FROM schedule_overrides
WHERE course_id = ?
  AND date = ?
  AND blocked = TRUE
`

type GetOverridesParams struct {
	CourseID int64
	Date     time.Time
}

type GetOverridesRow struct {
	StartTime time.Time
	EndTime   time.Time
	Blocked   bool
}

func (q *Queries) GetOverrides(ctx context.Context, arg *GetOverridesParams) ([]GetOverridesRow, error) {
	rows, err := q.db.QueryContext(ctx, getOverrides, arg.CourseID, arg.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOverridesRow
	for rows.Next() {
		var i GetOverridesRow
		if err := rows.Scan(&i.StartTime, &i.EndTime, &i.Blocked); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSchedule = `-- name: GetSchedule :one
SELECT course_id, TIMESTAMP(start_time) AS start_time, TIMESTAMP(end_time) AS end_time, occurrence, day, buffer
FROM schedule
WHERE course_id = ?
  AND day = ?
LIMIT 1
`

type GetScheduleParams struct {
	CourseID int64
	Day      int64
}

type GetScheduleRow struct {
	CourseID   int64
	StartTime  time.Time
	EndTime    time.Time
	Occurrence int64
	Day        int64
	Buffer     int64
}

func (q *Queries) GetSchedule(ctx context.Context, arg *GetScheduleParams) (GetScheduleRow, error) {
	row := q.db.QueryRowContext(ctx, getSchedule, arg.CourseID, arg.Day)
	var i GetScheduleRow
	err := row.Scan(
		&i.CourseID,
		&i.StartTime,
		&i.EndTime,
		&i.Occurrence,
		&i.Day,
		&i.Buffer,
	)
	return i, err
}

const getSchedules = `-- name: GetSchedules :many
SELECT course_id, TIMESTAMP(start_time) AS start_time, TIMESTAMP(end_time) AS end_time, occurrence, day, buffer
FROM schedule
WHERE course_id = ? ORDER BY day ASC
`

type GetSchedulesRow struct {
	CourseID   int64
	StartTime  time.Time
	EndTime    time.Time
	Occurrence int64
	Day        int64
	Buffer     int64
}

func (q *Queries) GetSchedules(ctx context.Context, courseID int64) ([]GetSchedulesRow, error) {
	rows, err := q.db.QueryContext(ctx, getSchedules, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSchedulesRow
	for rows.Next() {
		var i GetSchedulesRow
		if err := rows.Scan(
			&i.CourseID,
			&i.StartTime,
			&i.EndTime,
			&i.Occurrence,
			&i.Day,
			&i.Buffer,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertSchedule = `-- name: InsertSchedule :exec
INSERT INTO schedule(course_id, start_time, end_time, occurrence, day, buffer)
VALUES (?, ?, ?, ?, ?, ?)
`

type InsertScheduleParams struct {
	CourseID   int64
	StartTime  time.Time
	EndTime    time.Time
	Occurrence int64
	Day        int64
	Buffer     int64
}

func (q *Queries) InsertSchedule(ctx context.Context, arg *InsertScheduleParams) error {
	_, err := q.db.ExecContext(ctx, insertSchedule,
		arg.CourseID,
		arg.StartTime,
		arg.EndTime,
		arg.Occurrence,
		arg.Day,
		arg.Buffer,
	)
	return err
}
