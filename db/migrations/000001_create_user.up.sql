CREATE TABLE IF NOT EXISTS users (
    id varchar(36) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(255) NOT NULL,
    name varchar(100) NOT NULL,
    role varchar(5) NOT NULL,
    status boolean NOT NULL,
    last_login timestamp with time zone,
    reset_password_token varchar(255),
    verify_email_token varchar(255),
    is_verified boolean NOT NULL DEFAULT false,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
    );

ALTER TABLE users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE users
    ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'buyer'));

CREATE INDEX idx_users_deleted_at
    ON users USING btree
    (deleted_at ASC NULLS LAST);
