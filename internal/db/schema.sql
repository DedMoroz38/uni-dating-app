CREATE TABLE users (
  id   BIGSERIAL  NOT NULL PRIMARY KEY,
  username text    NOT NULL,
  email text    NOT NULL,
  password text    NOT NULL,
  created_at timestamp    NOT NULL,
  updated_at timestamp    NOT NULL
)