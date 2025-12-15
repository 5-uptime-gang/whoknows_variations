package main

import (
	"database/sql"
	"log"
	"time"
)

func InitDB(db *sql.DB) error {
	// 1) Create tables (PostgreSQL types + constraints)
	usersTable := `
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);`

	if _, err := db.Exec(usersTable); err != nil {
		return err
	}

	pagesTable := `
CREATE TABLE IF NOT EXISTS pages (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  language TEXT NOT NULL CHECK (language IN ('en', 'da')) DEFAULT 'en',
  last_updated TIMESTAMPTZ,
  content TEXT NOT NULL,
  tsv_document tsvector
);`

	if _, err := db.Exec(pagesTable); err != nil {
		return err
	}

	// 3) Enable search extensions, trigger, and indexes (idempotent)
	ftsSetup := `
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE OR REPLACE FUNCTION pages_tsvector_update() RETURNS trigger AS $$
BEGIN
  NEW.tsv_document :=
    setweight(
      to_tsvector(
        (CASE NEW.language WHEN 'da' THEN 'danish' ELSE 'english' END)::regconfig,
        coalesce(NEW.title, '')
      ),
      'A'
    )
    ||
    setweight(
      to_tsvector(
        (CASE NEW.language WHEN 'da' THEN 'danish' ELSE 'english' END)::regconfig,
        coalesce(NEW.content, '')
      ),
      'B'
    );
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS pages_tsvector_trigger ON pages;

CREATE TRIGGER pages_tsvector_trigger
BEFORE INSERT OR UPDATE ON pages
FOR EACH ROW EXECUTE FUNCTION pages_tsvector_update();

CREATE INDEX IF NOT EXISTS idx_pages_tsv_document
  ON pages USING GIN (tsv_document);

CREATE INDEX IF NOT EXISTS idx_pages_title_trgm
  ON pages USING GIN (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pages_content_trgm
  ON pages USING GIN (content gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_pages_last_updated
  ON pages (last_updated DESC);
`

	if _, err := db.Exec(ftsSetup); err != nil {
		return err
	}

	// 4) Seed admin user (SQLite: INSERT OR IGNORE -> PostgreSQL: ON CONFLICT DO NOTHING)
	seedAdmin := `
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
ON CONFLICT (username) DO NOTHING;`

	// NOTE: If you want to prevent duplicates by email as well, you can also choose:
	// ON CONFLICT DO NOTHING
	// but then you can't target a specific constraint/column set.
	if _, err := db.Exec(seedAdmin, "admin", "keamonk1@stud.kea.dk", "5f4dcc3b5aa765d61d8327deb882cf99"); err != nil {
		return err
	}

	// 3) Seed pages in a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// Safety rollback if Commit is not reached
		_ = tx.Rollback()
	}()

	insertPageStmt, err := tx.Prepare(`
INSERT INTO pages (title, url, language, last_updated, content)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (url) DO NOTHING;`)
	if err != nil {
		return err
	}
	defer func() {
		if err := insertPageStmt.Close(); err != nil {
			log.Printf("stmt.Close failed: %v", err)
		}
	}()

	now := time.Now()
	for _, page := range getPageSeedData() {
		if _, err := insertPageStmt.Exec(page.Title, page.URL, page.Language, now, page.Content); err != nil {
			log.Printf("Error inserting seed data (%s): %v", page.Title, err)
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getPageSeedData() []Page {
	return []Page{
		{"Go Basics", "https://go.dev/doc/tutorial/getting-started", "en", time.Time{}, "Go is a statically typed, compiled programming language designed at Google. Learn the basics of packages, functions, and goroutines."},
		{"SQL Joins", "https://www.w3schools.com/sql/sql_join.asp", "en", time.Time{}, "SQL joins are used to combine rows from two or more tables. Understand INNER JOIN, LEFT JOIN, RIGHT JOIN, and FULL OUTER JOIN."},
		{"Introduktion til Go", "https://go.dev/doc/", "da", time.Time{}, "Go (Golang) er et programmeringssprog udviklet af Google. Det er effektivt til backend-systemer og understøtter goroutines til samtidighed."},
		{"SQL Forespørgsler", "https://www.sqlitetutorial.net/", "da", time.Time{}, "SQL bruges til at hente og manipulere data i databaser. Eksempler inkluderer SELECT, INSERT, UPDATE og DELETE forespørgsler."},
		{"Python Basics", "https://docs.python.org/3/tutorial/", "en", time.Time{}, "Python is a high-level, interpreted language with dynamic typing. Learn about variables, functions, and control flow."},
		{"JavaScript Guide", "https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide", "en", time.Time{}, "JavaScript is the programming language of the web. Learn core concepts like functions, objects, and asynchronous programming."},
		{"HTML Introduktion", "https://www.w3schools.com/html/", "da", time.Time{}, "HTML er struktursproget for internettet. Lær om tags, elementer, links og tabeller."},
		{"CSS Layout", "https://developer.mozilla.org/en-US/docs/Learn/CSS/CSS_layout", "en", time.Time{}, "CSS is used to style web pages. Learn about Flexbox, Grid, and positioning."},
		{"C# Programming Guide", "https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/", "en", time.Time{}, "C# is a modern, object-oriented programming language. Learn about classes, interfaces, and LINQ queries."},
		{"Java Grundlæggende", "https://www.javatpoint.com/java-tutorial", "da", time.Time{}, "Java er et populært objektorienteret sprog. Det bruges til webapps, mobilapps og enterprise-løsninger."},
		{"Rust Language Book", "https://doc.rust-lang.org/book/", "en", time.Time{}, "Rust is a systems programming language focusing on safety and performance. Learn about ownership, borrowing, and concurrency."},
		{"Kotlin Basics", "https://kotlinlang.org/docs/basic-syntax.html", "en", time.Time{}, "Kotlin is a modern language for JVM and Android development. Learn about null safety, data classes, and coroutines."},
		{"Linux Kommandolinje", "https://linuxcommand.org/lc3_learning_the_shell.php", "da", time.Time{}, "Linux shell giver dig kontrol over systemet. Lær kommandoer som ls, cd, grep og pipes."},
		{"Git Tutorial", "https://git-scm.com/docs/gittutorial", "en", time.Time{}, "Git is a distributed version control system. Learn about commits, branches, merging, and rebasing."},
		{"Docker Basics", "https://docs.docker.com/get-started/", "en", time.Time{}, "Docker allows you to package applications into containers. Learn about images, containers, and Docker Compose."},
		{"Introduktion til Kubernetes", "https://kubernetes.io/da/docs/tutorials/kubernetes-basics/", "da", time.Time{}, "Kubernetes automatiserer udrulning og skalering af containeriserede applikationer."},
		{"React Documentation", "https://react.dev/learn", "en", time.Time{}, "React is a JavaScript library for building user interfaces. Learn about components, props, and hooks."},
		{"Node.js Guide", "https://nodejs.org/en/docs/", "en", time.Time{}, "Node.js is a runtime environment for executing JavaScript outside the browser. Learn about modules, streams, and async I/O."},
		{"TypeScript Handbook", "https://www.typescriptlang.org/docs/", "en", time.Time{}, "TypeScript is a superset of JavaScript that adds static typing. Learn about interfaces, generics, and decorators."},
		{"Machine Learning Basics", "https://scikit-learn.org/stable/tutorial/basic/tutorial.html", "en", time.Time{}, "Machine Learning is about teaching computers to learn patterns from data. Learn about supervised and unsupervised learning."},
		{"Python Data Analysis", "https://pandas.pydata.org/docs/getting_started/index.html", "en", time.Time{}, "Pandas is a Python library for data analysis. Learn about DataFrames, indexing, and data cleaning."},
		{"Introduktion til MySQL", "https://www.mysqltutorial.org/", "da", time.Time{}, "MySQL er en open source relationsdatabase. Lær hvordan man opretter tabeller, indsætter data og udfører forespørgsler."},
		{"Cybersecurity Fundamentals", "https://www.coursera.org/learn/cyber-security-basics", "en", time.Time{}, "Cybersecurity focuses on protecting systems and data from attacks. Learn about encryption, firewalls, and threat modeling."},
		{"AI og Maskinlæring", "https://da.wikipedia.org/wiki/Maskinl%C3%A6ring", "da", time.Time{}, "Maskinlæring er en gren af kunstig intelligens, hvor algoritmer trænes til at finde mønstre i data."},
	}
}
