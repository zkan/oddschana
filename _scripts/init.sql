CREATE TABLE people
(
    mobile_number TEXT NOT NULL
);

CREATE TABLE places
(
    name TEXT NOT NULL,
    limits INT NOT NULL,
    lat REAL,
    long REAL
);

CREATE TABLE visit
(
    people_id INTEGER NOT NULL,
    place_id INTEGER NOT NULL
);
