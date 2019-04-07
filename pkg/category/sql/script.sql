CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE categories (
    uid UUID PRIMARY KEY,
    user_uid UUID NOT NULL,
    name VARCHAR(80) NOT NULL
);
