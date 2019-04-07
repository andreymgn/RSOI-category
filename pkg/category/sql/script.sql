CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE categories (
    uid UUID PRIMARY KEY,
    user_uid UUID NOT NULL,
    name VARCHAR(80) NOT NULL
);

CREATE TABLE reports (
    uid UUID PRIMARY KEY,
    category_uid UUID NOT NULL,
    post_uid UUID NOT NULL,
    comment_uid UUID,
    reason VARCHAR(160) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);