CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email TEXT NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       role TEXT NOT NULL CHECK (role IN ('employee', 'moderator'))
);

CREATE TABLE pvz (
                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                     city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')),
                     registration_date TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE reception (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
                           date_time TIMESTAMP NOT NULL DEFAULT NOW(),
                           status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE product (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         reception_id UUID NOT NULL REFERENCES reception(id) ON DELETE CASCADE,
                         type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
                         date_time TIMESTAMP NOT NULL DEFAULT NOW()
);
