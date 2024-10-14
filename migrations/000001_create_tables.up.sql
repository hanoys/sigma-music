CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    start_date TIMESTAMP NOT NULL,
    expiration_date TIMESTAMP NOT NULL CHECK (expiration_date > start_date)
);

CREATE TABLE IF NOT EXISTS musicians (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    description VARCHAR(1024) NOT NULL
);

CREATE TABLE IF NOT EXISTS albums (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    published BOOLEAN NOT NULL,
    release_date TIMESTAMP
);

CREATE TABLE IF NOT EXISTS album_musician (
    musician_id UUID REFERENCES musicians ON DELETE CASCADE,
    album_id UUID REFERENCES albums ON DELETE CASCADE,
    PRIMARY KEY (musician_id, album_id)
);

CREATE TABLE IF NOT EXISTS tracks (
    id UUID PRIMARY KEY,
    album_id UUID REFERENCES albums ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048) NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
    id UUID PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS track_genre (
    track_id UUID REFERENCES tracks ON DELETE CASCADE,
    genre_id UUID REFERENCES genres ON DELETE CASCADE,
    PRIMARY KEY (track_id, genre_id)
);

CREATE TABLE IF NOT EXISTS favorite (
    user_id UUID REFERENCES users ON DELETE CASCADE,
    track_id UUID REFERENCES tracks ON DELETE CASCADE,
    PRIMARY KEY (user_id, track_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    track_id UUID REFERENCES tracks ON DELETE CASCADE,
    stars INTEGER CHECK (stars > 0 AND stars <= 5),
    comment_text VARCHAR(1024)
);

CREATE TABLE IF NOT EXISTS users_history (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users ON DELETE CASCADE,
    track_id UUID REFERENCES tracks ON DELETE CASCADE
);
