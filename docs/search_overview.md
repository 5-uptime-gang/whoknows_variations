# Search & Ingest Overview

This document explains how search works across the Go API and the ingest Node.js helper, along with the database structures that support them.

## Database search structures
- Extensions: `pg_trgm` enabled for trigram similarity/ILIKE fallback.
- FTS materialization: `pages.tsv_document` maintained by trigger `pages_tsvector_update` (not a stored generated column). It weights title higher than content and uses `danish`/`english` regconfig per row language.
- Indexes:
  - `GIN` on `tsv_document` for fast FTS.
  - `GIN` trigram indexes on `title` and `content` for typo/partial matching.
  - `last_updated` index for recency tie-break.
- Primary key: `id BIGSERIAL` is the table PK; `url` remains unique; `title` can repeat.
- Migrations:
  - `migrations/001_full_text_search.sql` sets up extensions, tsvector, and indexes.
  - `migrations/002_pages_id_pk.sql` moves the primary key to `id BIGSERIAL` (URL stays unique; titles may repeat).
  - `InitDB` mirrors this combined setup so fresh databases match the migrations.

## Go API search flow
- Endpoint: `GET /api/search?q=...&language=...&limit=...`
- Language detection: respects `language` query param if provided (`en`/`da`), otherwise heuristics (Danish characters/stopwords) to choose `danish` vs `english`.
- Query plan:
  1) Full-text search on `tsv_document` with `plainto_tsquery`, boosting title via weights and adding a small recency score.
  2) If FTS doesn’t fill the requested limit, fallback runs trigram + `ILIKE` over title/content with title weighted higher.
  3) Results include `title`, `url`, `language`, `last_updated`, and a `snippet` via `ts_headline`.
- Safety/perf: parameterized queries, capped limit (1–50), uses indexes above.

## Node.js ingest/search helpers (`search-ingest`)
- Ingest runner: `npm run ingest` scrapes/clusters queries and writes `pages.json`.
- Import into Postgres: `npm run import -- --database-url postgres://user:pass@host:port/db --file pages.json` upserts pages into the `pages` table (defaults to `DATABASE_URL` and `pages.json`).
- Search helper: `src/search/searchPages.js` exposes `searchPages(query, { limit, language })` and `detectLanguage`, mirroring the Go API logic (FTS first, trigram fallback, snippets).

## Frontend display
- Search results use `page.snippet` (with `<b>` highlights) and clamp descriptions to 5 lines via CSS (`.search-result-description` uses `-webkit-line-clamp: 5` with ellipsis).

## How to apply and test
1) Apply migration: `psql "$DATABASE_URL" -f migrations/001_full_text_search.sql`.
2) Import pages: `cd search-ingest && npm run import -- --database-url "$DATABASE_URL"`.
3) Exercise API: `curl "http://localhost:8080/api/search?q=go routines&limit=5"` or Danish example `curl "http://localhost:8080/api/search?q=hvordan skriver jeg go kode&language=da"`.
