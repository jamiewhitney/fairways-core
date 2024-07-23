CREATE UNIQUE INDEX idx_course_datetime ON bookings(course_id, datetime);
# DROP INDEX idx_course_datetime ON bookings;