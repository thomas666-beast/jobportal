-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION immutable_array_to_string(arr text[], sep text)
RETURNS text
LANGUAGE sql
IMMUTABLE
AS $$
    SELECT array_to_string(arr, sep);
$$;
-- +goose StatementEnd

ALTER TABLE jobs
    ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(immutable_array_to_string(tags, ' '), '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(description, '')), 'C')
    ) STORED;

CREATE INDEX idx_jobs_search_vector ON jobs USING GIN (search_vector);

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_jobs_search_vector;
ALTER TABLE jobs DROP COLUMN IF EXISTS search_vector;
DROP FUNCTION IF EXISTS immutable_array_to_string(text[], text);
-- +goose StatementEnd
