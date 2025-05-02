CREATE TABLE IF NOT EXISTS outbox (
    "id" VARCHAR(100) NOT NULL,
    "data" BLOB NOT NULL,
    "state" INT NOT NULL,
    "created_on" TIMESTAMP NOT NULL,
    "locked_by" VARCHAR(100) NULL,
    "locked_on" TIMESTAMP NULL,
    "processed_on" TIMESTAMP NULL,
    "number_of_attempts" INT NOT NULL,
    "last_attempted_on" TIMESTAMP NULL,
    "error" VARCHAR(1000) NULL
)