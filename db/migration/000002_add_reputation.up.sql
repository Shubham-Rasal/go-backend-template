-- update the schema to add the reputation column
ALTER TABLE users ADD COLUMN reputation integer NOT NULL DEFAULT 0;