CREATE TABLE IF NOT EXISTS users (
    id PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
);

CREATE TABLE IF NOT EXISTS preferences (
    id PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    units TEXT NOT NULL CHECK (units IN ('us', 'si')) DEFAULT 'us',
    temperature_min REAL DEFAULT 23.9, -- Minimum temperature in Celsius (~75°F)
    temperature_max REAL DEFAULT 38.0, -- Maximum temperature in Celsius (~100°F)
    wind_min INT DEFAULT 0, -- beaufort scale minimum
    wind_max INT DEFAULT 4, -- beaufort scale maximum
    precipitation_min REAL DEFAULT 0.0, -- Minimum precipitation prob
    precipitation_max REAL DEFAULT 10.0, -- Maximum precipitation prob
);

CREATE TABLE IF NOT EXISTS likes (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    liked INT NOT NULL,
    temperature REAL NOT NULL,
    wind INT NOT NULL,
    precipitation REAL NOT NULL,
    PRIMARY KEY (user_id, temperature, wind, precipitation)
);