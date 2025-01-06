CREATE TABLE posts_new (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    user_id bigint NOT NULL,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    tags varchar(100) [],
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

INSERT INTO posts_new (id, title, user_id, content, created_at, tags, updated_at)
SELECT id, title, user_id, content, created_at, tags, updated_at
FROM posts;

ALTER TABLE posts_new
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);

DROP TABLE posts;

ALTER TABLE posts_new RENAME TO posts;