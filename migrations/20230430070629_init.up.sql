CREATE TABLE IF NOT EXISTS persons (
    uid Integer NOT NULL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL
);

CREATE INDEX idx_persons_first_name ON persons (first_name);
CREATE INDEX idx_persons_last_name ON persons (last_name);
COMMENT ON TABLE persons IS 'Таблица имен людей';
COMMENT ON COLUMN persons.uid IS 'uid from xml';
COMMENT ON COLUMN persons.first_name IS 'first name';
COMMENT ON COLUMN persons.last_name IS 'last name';