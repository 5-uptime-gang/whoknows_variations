-- Enable trigram matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add tsvector column (regular column, not generated)
ALTER TABLE pages
ADD COLUMN IF NOT EXISTS tsv_document tsvector;

-- Create trigger function to update tsv_document
CREATE OR REPLACE FUNCTION pages_tsvector_update() RETURNS trigger AS $$
BEGIN
  NEW.tsv_document :=
    setweight(
      to_tsvector(
        CASE NEW.language WHEN 'da' THEN 'danish' ELSE 'english' END::regconfig,
        coalesce(NEW.title, '')
      ),
      'A'
    )
    ||
    setweight(
      to_tsvector(
        CASE NEW.language WHEN 'da' THEN 'danish' ELSE 'english' END::regconfig,
        coalesce(NEW.content, '')
      ),
      'B'
    );

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger
DROP TRIGGER IF EXISTS pages_tsvector_trigger ON pages;

CREATE TRIGGER pages_tsvector_trigger
BEFORE INSERT OR UPDATE ON pages
FOR EACH ROW EXECUTE FUNCTION pages_tsvector_update();

-- Backfill existing rows
UPDATE pages SET tsv_document = NULL;

-- Indexes
CREATE INDEX IF NOT EXISTS idx_pages_tsv_document
ON pages USING GIN (tsv_document);

CREATE INDEX IF NOT EXISTS idx_pages_title_trgm
ON pages USING GIN (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pages_content_trgm
ON pages USING GIN (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pages_last_updated
ON pages (last_updated DESC);
