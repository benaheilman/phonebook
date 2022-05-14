DROP INDEX IF EXISTS listing_phone;
DROP TABLE IF EXISTS listing;
CREATE TABLE listing (
    id INTEGER PRIMARY KEY,
    name TEXT,
    surname TEXT NOT NULL,
    phone TEXT NOT NULL,
    updated TIMESTAMP
);
CREATE UNIQUE INDEX listing_phone ON listing (phone);