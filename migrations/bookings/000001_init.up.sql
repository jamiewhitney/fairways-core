CREATE DATABASE IF NOT EXISTS bookings;
CREATE TABLE IF NOT EXISTS bookings
(
    id                BIGINT PRIMARY KEY AUTO_INCREMENT                               NOT NULL,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP                             NOT NULL,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL,
    deleted_at        TIMESTAMP                                                       NULL,
    user_id           VARCHAR(255)                                                    NOT NULL,
    course_id         BIGINT                                                          NOT NULL,
    golfers           BIGINT                                                          NOT NULL,
    datetime          TIMESTAMP                                                       NOT NULL,
    price             FLOAT                                                           NOT NULL,
    booking_id        VARCHAR(32)                                                     NOT NULL,
    stripe_payment_id VARCHAR(32)                                                     NOT NULL,
    status            ENUM ('confirmed', 'requested')                                 NOT NULL DEFAULT 'requested',
    confirmed         BOOL                                                            NOT NULL DEFAULT FALSE
);