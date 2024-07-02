CREATE TABLE IF NOT EXISTS candidates (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name TEXT not null,
  created_at TIMESTAMP default now()
  );