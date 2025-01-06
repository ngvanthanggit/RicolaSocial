CREATE TABLE posts_new (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    tags varchar(100) [],
    user_id bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

INSERT INTO posts_new (id, title, content, tags, user_id, created_at, updated_at)
SELECT id, title, content, tags, user_id, created_at, updated_at
FROM posts;

ALTER TABLE posts_new
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id);

DROP TABLE posts;

ALTER TABLE posts_new RENAME TO posts;