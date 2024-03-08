ALTER TABLE episodes ADD CONSTRAINT movies_runtime_check CHECK (runtime >= 0);
ALTER TABLE episodes ADD CONSTRAINT movies_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE episodes ADD CONSTRAINT genres_length_check CHECK (array_length(characters, 1) BETWEEN 1 AND 5);