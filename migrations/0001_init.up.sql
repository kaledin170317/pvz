-- Включаем расширение UUID
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pvz
(
    id                UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city              TEXT      NOT NULL CHECK (city IN ('Санкт-Петербург', 'Казань', 'Москва'))
);

CREATE TABLE reception
(
    id        UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    date_time TIMESTAMP NOT NULL DEFAULT NOW(),
    pvz_id    UUID      NOT NULL REFERENCES pvz (id),
    status    TEXT      NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE product
(
    id           UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    date_time    TIMESTAMP NOT NULL DEFAULT NOW(),
    type         TEXT      NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID      NOT NULL REFERENCES reception (id)
);

CREATE TABLE users
(
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role     TEXT NOT NULL CHECK (role IN ('employee', 'moderator', 'client'))
);
