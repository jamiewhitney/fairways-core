CREATE DATABASE IF NOT EXISTS pricing;

CREATE TABLE IF NOT EXISTS course_base_prices(id INT PRIMARY KEY AUTO_INCREMENT,course_id INT NOT NULL,price INT NOT NULL, year INT NOT NULL);

CREATE TABLE IF NOT EXISTS course_pricing_rules(
    id                 INT PRIMARY KEY AUTO_INCREMENT,
    price_band         INT       NOT NULL,
    week               INT       NOT NULL,
    hour               INT       NOT NULL,
    golfers_1_modifier FLOAT     NOT NULL DEFAULT 1,
    golfers_2_modifier FLOAT     NOT NULL DEFAULT 0.85,
    golfers_3_modifier FLOAT     NOT NULL DEFAULT 0.75,
    golfers_4_modifier FLOAT     NOT NULL DEFAULT 0.65,
    created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS course_price_bands(
    id        INT PRIMARY KEY AUTO_INCREMENT,
    band      INT NOT NULL,
    course_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS missed_price_requests(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id  INT       NOT NULL,
    year       INT       NOT NULL,
    week       INT       NOT NULL,
    hour       INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
