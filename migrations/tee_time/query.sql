-- name: GetSchedule :one
SELECT course_id, TIMESTAMP(start_time) AS start_time, TIMESTAMP(end_time) AS end_time, occurrence, day, buffer
FROM schedule
WHERE course_id = ?
  AND day = ?
LIMIT 1;
-- name: GetOverrides :many
SELECT start_time, end_time, blocked
FROM schedule_overrides
WHERE course_id = ?
  AND date = ?
  AND blocked = TRUE;

-- name: GetSchedules :many
SELECT course_id, TIMESTAMP(start_time) AS start_time, TIMESTAMP(end_time) AS end_time, occurrence, day, buffer
FROM schedule
WHERE course_id = ? ORDER BY day ASC;


-- name: InsertSchedule :exec
INSERT INTO schedule(course_id, start_time, end_time, occurrence, day, buffer)
VALUES (?, ?, ?, ?, ?, ?);