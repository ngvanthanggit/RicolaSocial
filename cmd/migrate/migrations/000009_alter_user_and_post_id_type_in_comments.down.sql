CREATE SEQUENCE IF NOT EXISTS comments_user_id_seq;

ALTER TABLE comments ALTER COLUMN user_id
SET DEFAULT nextval('comments_user_id_seq'::regclass);

SELECT setval('comments_user_id_seq', COALESCE((SELECT MAX(user_id) FROM comments), 1));


CREATE SEQUENCE IF NOT EXISTS comments_post_id_seq;

ALTER TABLE comments ALTER COLUMN post_id
SET DEFAULT nextval('comments_post_id_seq'::regclass);

SELECT setval('comments_post_id_seq', COALESCE((SELECT MAX(post_id) FROM comments), 1));
