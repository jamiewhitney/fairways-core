CREATE DATABASE IF NOT EXISTS tee_times;
CREATE TABLE IF NOT EXISTS schedule
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    course_id  BIGINT NOT NULL,
    start_time TIME   NOT NULL,
    end_time   TIME   NOT NULL,
    occurrence BIGINT NOT NULL,
    day        BIGINT NOT NULL,
    UNIQUE KEY unique_course_day (course_id, day)
);

CREATE TABLE IF NOT EXISTS schedule_overrides
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME  DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    course_id  BIGINT    NOT NULL,
    start_time TIME NOT NULL,
    end_time   TIME NOT NULL,
    occurrence BIGINT    NOT NULL,
    date       DATE      NOT NULL,
    blocked    BOOL      NOT NULL
);