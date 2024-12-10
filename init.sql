CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    avatar TEXT,
    comment TEXT
    );