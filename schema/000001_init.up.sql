CREATE TABLE users
(
    id serial not NULL UNIQUE PRIMARY KEY,
    balance double PRECISION CONSTRAINT positive_balance CHECK (balance >= 0)

);

CREATE TABLE operations
(
    id serial not NULL UNIQUE,
    main_subject_id int,
    other_subject_id int,
    amount_of_money DOUBLE PRECISION,
    executed_at timestamp,
    PRIMARY KEY (id),
    CONSTRAINT FK_users_operations 
        FOREIGN KEY(main_subject_id)
            REFERENCES users(id),
        FOREIGN KEY (other_subject_id)
            REFERENCES users(id)
)