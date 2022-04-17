CREATE TABLE IF NOT EXISTS ApplicationUser (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS UserFeed (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS Feed (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feedName VARCHAR NOT NULL,
    url VARCHAR NOT NULL
);