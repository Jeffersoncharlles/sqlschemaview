CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title TEXT,
    CONSTRAINT fk_posts_users
        FOREIGN KEY (user_id)
        REFERENCES users(id)
);
