CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO users (username, email, password)
VALUES ('thangitcbg', 
        'thangitcbg@gmail.com', 
        convert_to(crypt('atdznvn2003', gen_salt('bf')), 'UTF8')
);