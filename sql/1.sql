--! This is the SQL script for initializing the TU Password Research database !--

--! TABLES !--
CREATE TABLE buildings (
    id                  BIGSERIAL               PRIMARY KEY,
    uid                 TEXT,
    name                TEXT                    NOT NULL,
    description         TEXT,
    image_ref           TEXT,
    longitude           DOUBLE PRECISION        NOT NULL,
    latitude            DOUBLE PRECISION,
    address             TEXT
);

CREATE TABLE rooms (
    id                  BIGSERIAL               PRIMARY KEY,
    building_id         INTEGER                 REFERENCES buildings(id),
    floor               INTEGER,
    room_number         VARCHAR(10),
    details             TEXT
);

CREATE table user_search_query (
    id                  BIGSERIAL               PRIMARY KEY,
    query               TEXT,
    building_id         INTEGER                 REFERENCES buildings(id),
    creation_date       TIMESTAMP
);

CREATE TABLE version_number (
    id                  BIGSERIAL               PRIMARY KEY,
    version             DOUBLE PRECISION
)

CREATE TABLE images
(
    id                  BIGSERIAL               PRIMARY KEY,
    name                TEXT                    NOT NULL,
    zoom                INTEGER                 NOT NULL,
    lat                 DOUBLE PRECISION        NOT NULL,
    "long"              DOUBLE PRECISION        NOT NULL,
    width               INTEGER                 NOT NULL,
    height              INTEGER                 NOT NULL,
    image               TEXT                    NOT NULL
)

CREATE TABLE labels
(
    id                  BIGSERIAL               PRIMARY KEY,
    name                TEXT                    NOT NULL,
    zoom                INTEGER                 NOT NULL,
    color               INTEGER                 NOT NULL,
    size                INTEGER                 NOT NULL,
    lat                 DOUBLE PRECISION        NOT NULL,
    "long"              DOUBLE PRECISION        NOT NULL
)
