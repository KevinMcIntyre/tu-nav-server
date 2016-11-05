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
