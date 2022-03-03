CREATE TABLE users
(
    id serial not NULL UNIQUE PRIMARY KEY,
    balance double PRECISION CONSTRAINT positive_balance CHECK (balance >= 0)

);

CREATE TABLE operations
(
    id serial not NULL UNIQUE
    user_id
)