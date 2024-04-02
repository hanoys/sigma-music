CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users ON DELETE CASCADE,
    create_time TIMESTAMP NOT NULL,
    price NUMERIC(20, 5) NOT NULL CHECK (price > 0)
);

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users ON DELETE CASCADE,
    order_id INTEGER REFERENCES orders ON DELETE CASCADE,
    start_date TIMESTAMP NOT NULL,
    expiration_date TIMESTAMP NOT NULL CHECK (expiration_date > start_date)
);

CREATE TABLE IF NOT EXISTS musicians (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    description VARCHAR(1024) NOT NULL
);

CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1024),
    release_date TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS album_musician (
    musician_id INTEGER REFERENCES musicians ON DELETE CASCADE,
    album_id INTEGER REFERENCES albums ON DELETE CASCADE,
    PRIMARY KEY (musician_id, album_id)
);

CREATE TABLE IF NOT EXISTS tracks (
    id SERIAL PRIMARY KEY,
    album_id INTEGER REFERENCES albums ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2048) NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS track_genre (
    id SERIAL PRIMARY KEY,
    track_id INTEGER REFERENCES tracks ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genres ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS favorite (
    user_id INTEGER REFERENCES users ON DELETE CASCADE,
    track_id INTEGER REFERENCES tracks ON DELETE CASCADE,
    PRIMARY KEY (user_id, track_id)
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users ON DELETE CASCADE,
    track_id INTEGER REFERENCES tracks ON DELETE CASCADE,
    stars INTEGER CHECK (stars > 0 AND stars <= 5),
    comment_text VARCHAR(1024)
);
