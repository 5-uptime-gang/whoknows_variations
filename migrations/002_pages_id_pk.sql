-- Migrate pages table to use an integer primary key instead of title

-- Drop old primary key on title (if present)
ALTER TABLE pages DROP CONSTRAINT IF EXISTS pages_pkey;

-- Add id column if missing
ALTER TABLE pages
    ADD COLUMN IF NOT EXISTS id BIGSERIAL;

-- Ensure existing rows get an id
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'pages' AND column_name = 'id') THEN
        UPDATE pages
        SET id = nextval(pg_get_serial_sequence('pages', 'id'))
        WHERE id IS NULL;
    END IF;
END$$;

-- Ensure id is primary key
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'pages'::regclass
          AND contype = 'p'
    ) THEN
        ALTER TABLE pages ADD CONSTRAINT pages_pkey PRIMARY KEY (id);
    END IF;
END$$;

-- Title should no longer be unique; uniqueness remains on URL
