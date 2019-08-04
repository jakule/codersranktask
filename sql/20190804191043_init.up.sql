CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE secrets
(
    id     uuid DEFAULT uuid_generate_v4(),
    secret VARCHAR NOT NULL,
    PRIMARY KEY (id)
);