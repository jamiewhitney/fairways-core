-- name: GetABooking :one
SELECT *
FROM bookings
WHERE id = ?
LIMIT 1;

-- name: ListUserBookings :many
SELECT *
FROM bookings
WHERE user_id = ?
ORDER BY datetime;

-- name: CreateBooking :execlastid
INSERT INTO bookings(user_id, course_id, golfers, datetime, price, booking_id, stripe_payment_id)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateBooking :exec
UPDATE bookings
SET user_id = ?
WHERE id = ?;

-- name: UpdateAndConfirm :exec
UPDATE bookings
SET status            = 'confirmed',
    confirmed         = 1,
    stripe_payment_id = ?
WHERE id = ?;

-- name: GetConfirmedBookingsByDateAndCourse :many
SELECT *
FROM bookings
WHERE status = 'confirmed'
  AND course_id = ?
  AND DATE(datetime) = ?;