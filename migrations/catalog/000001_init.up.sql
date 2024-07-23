CREATE TABLE IF NOT EXISTS courses
(
    id               BIGINT AUTO_INCREMENT PRIMARY KEY,
    name             VARCHAR(255)                         NOT NULL,
    holes            INT                                  NOT NULL,
    lapsed           TINYINT(1) DEFAULT 0                 NOT NULL,
    live             TINYINT(1) DEFAULT 0                 NOT NULL,
    street_address_1 VARCHAR(255)                         NOT NULL,
    city             VARCHAR(255)                         NOT NULL,
    state            VARCHAR(255)                         NOT NULL,
    postal_code      VARCHAR(255)                         NOT NULL,
    country          VARCHAR(255)                         NOT NULL,
    latitude         DECIMAL(10, 6)                       NULL,
    longitude        DECIMAL(10, 6)                       NULL,
    created_at       TIMESTAMP  DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at       TIMESTAMP  DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS holes
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    course_id  BIGINT                              NOT NULL,
    name       VARCHAR(255)                        NOT NULL,
    number     INT                                 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT holes_ibfk_1
        FOREIGN KEY (course_id) REFERENCES courses (id)
);

CREATE INDEX course_id
    ON holes (course_id);

CREATE TABLE IF NOT EXISTS tee_box
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    hole_id    BIGINT                              NOT NULL,
    color      VARCHAR(255)                        NOT NULL,
    yardage    INT                                 NOT NULL,
    rating     FLOAT                               NOT NULL,
    slope      INT                                 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT tee_box_ibfk_1
        FOREIGN KEY (hole_id) REFERENCES courses (id)
);

CREATE INDEX hole_id
    ON tee_box (hole_id);