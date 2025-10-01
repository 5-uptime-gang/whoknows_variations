package main

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func InitDB(db *sql.DB) error {
	// Users: create + drop - dropper users og laver den efterfølgende
	usersSchema := `
	DROP TABLE IF EXISTS users;

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);

	INSERT OR IGNORE INTO users (username, email, password) 
	VALUES ('admin', 'keamonk1@stud.kea.dk', '5f4dcc3b5aa765d61d8327deb882cf99');`

	if _, err := db.Exec(usersSchema); err != nil {
		return err
	}

	// Pages: create + drop
	pagesSchema := `
	CREATE TABLE IF NOT EXISTS pages (
		title TEXT PRIMARY KEY UNIQUE,
		url TEXT NOT NULL UNIQUE,
		language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
		last_updated TIMESTAMP,
		content TEXT NOT NULL
	);`

	if _, err := db.Exec(pagesSchema); err != nil {
		return err
	}

	// bruger prepared statement til at indsætte dataen
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO pages (title, url, language, last_updated, content) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	seedData := []struct {
		Title   string
		URL     string
		Lang    string
		Content string
	}{
		{
			Title:   "Go Basics",
			URL:     "https://go.dev/doc/tutorial/getting-started",
			Lang:    "en",
			Content: "Go is a statically typed, compiled programming language designed at Google. Learn the basics of packages, functions, and goroutines.",
		},
		{
			Title:   "SQL Joins",
			URL:     "https://www.w3schools.com/sql/sql_join.asp",
			Lang:    "en",
			Content: "SQL joins are used to combine rows from two or more tables. Understand INNER JOIN, LEFT JOIN, RIGHT JOIN, and FULL OUTER JOIN.",
		},
		{
			Title:   "Introduktion til Go",
			URL:     "https://go.dev/doc/",
			Lang:    "da",
			Content: "Go (Golang) er et programmeringssprog udviklet af Google. Det er effektivt til backend-systemer og understøtter goroutines til samtidighed.",
		},
		{
			Title:   "SQL Forespørgsler",
			URL:     "https://www.sqlitetutorial.net/",
			Lang:    "da",
			Content: "SQL bruges til at hente og manipulere data i databaser. Eksempler inkluderer SELECT, INSERT, UPDATE og DELETE forespørgsler.",
		},
				{
			Title:   "Python Basics",
			URL:     "https://docs.python.org/3/tutorial/",
			Lang:    "en",
			Content: "Python is a high-level, interpreted language with dynamic typing. Learn about variables, functions, and control flow.",
		},
		{
			Title:   "JavaScript Guide",
			URL:     "https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide",
			Lang:    "en",
			Content: "JavaScript is the programming language of the web. Learn core concepts like functions, objects, and asynchronous programming.",
		},
		{
			Title:   "HTML Introduktion",
			URL:     "https://www.w3schools.com/html/",
			Lang:    "da",
			Content: "HTML er struktursproget for internettet. Lær om tags, elementer, links og tabeller.",
		},
		{
			Title:   "CSS Layout",
			URL:     "https://developer.mozilla.org/en-US/docs/Learn/CSS/CSS_layout",
			Lang:    "en",
			Content: "CSS is used to style web pages. Learn about Flexbox, Grid, and positioning.",
		},
		{
			Title:   "C# Programming Guide",
			URL:     "https://learn.microsoft.com/en-us/dotnet/csharp/programming-guide/",
			Lang:    "en",
			Content: "C# is a modern, object-oriented programming language. Learn about classes, interfaces, and LINQ queries.",
		},
		{
			Title:   "Java Grundlæggende",
			URL:     "https://www.javatpoint.com/java-tutorial",
			Lang:    "da",
			Content: "Java er et populært objektorienteret sprog. Det bruges til webapps, mobilapps og enterprise-løsninger.",
		},
		{
			Title:   "Rust Language Book",
			URL:     "https://doc.rust-lang.org/book/",
			Lang:    "en",
			Content: "Rust is a systems programming language focusing on safety and performance. Learn about ownership, borrowing, and concurrency.",
		},
		{
			Title:   "Kotlin Basics",
			URL:     "https://kotlinlang.org/docs/basic-syntax.html",
			Lang:    "en",
			Content: "Kotlin is a modern language for JVM and Android development. Learn about null safety, data classes, and coroutines.",
		},
		{
			Title:   "Linux Kommandolinje",
			URL:     "https://linuxcommand.org/lc3_learning_the_shell.php",
			Lang:    "da",
			Content: "Linux shell giver dig kontrol over systemet. Lær kommandoer som ls, cd, grep og pipes.",
		},
		{
			Title:   "Git Tutorial",
			URL:     "https://git-scm.com/docs/gittutorial",
			Lang:    "en",
			Content: "Git is a distributed version control system. Learn about commits, branches, merging, and rebasing.",
		},
		{
			Title:   "Docker Basics",
			URL:     "https://docs.docker.com/get-started/",
			Lang:    "en",
			Content: "Docker allows you to package applications into containers. Learn about images, containers, and Docker Compose.",
		},
		{
			Title:   "Introduktion til Kubernetes",
			URL:     "https://kubernetes.io/da/docs/tutorials/kubernetes-basics/",
			Lang:    "da",
			Content: "Kubernetes automatiserer udrulning og skalering af containeriserede applikationer.",
		},
		{
			Title:   "React Documentation",
			URL:     "https://react.dev/learn",
			Lang:    "en",
			Content: "React is a JavaScript library for building user interfaces. Learn about components, props, and hooks.",
		},
		{
			Title:   "Node.js Guide",
			URL:     "https://nodejs.org/en/docs/",
			Lang:    "en",
			Content: "Node.js is a runtime environment for executing JavaScript outside the browser. Learn about modules, streams, and async I/O.",
		},
		{
			Title:   "TypeScript Handbook",
			URL:     "https://www.typescriptlang.org/docs/",
			Lang:    "en",
			Content: "TypeScript is a superset of JavaScript that adds static typing. Learn about interfaces, generics, and decorators.",
		},
		{
			Title:   "Machine Learning Basics",
			URL:     "https://scikit-learn.org/stable/tutorial/basic/tutorial.html",
			Lang:    "en",
			Content: "Machine Learning is about teaching computers to learn patterns from data. Learn about supervised and unsupervised learning.",
		},
		{
			Title:   "Python Data Analysis",
			URL:     "https://pandas.pydata.org/docs/getting_started/index.html",
			Lang:    "en",
			Content: "Pandas is a Python library for data analysis. Learn about DataFrames, indexing, and data cleaning.",
		},
		{
			Title:   "Introduktion til MySQL",
			URL:     "https://www.mysqltutorial.org/",
			Lang:    "da",
			Content: "MySQL er en open source relationsdatabase. Lær hvordan man opretter tabeller, indsætter data og udfører forespørgsler.",
		},
		{
			Title:   "Cybersecurity Fundamentals",
			URL:     "https://www.coursera.org/learn/cyber-security-basics",
			Lang:    "en",
			Content: "Cybersecurity focuses on protecting systems and data from attacks. Learn about encryption, firewalls, and threat modeling.",
		},
		{
			Title:   "AI og Maskinlæring",
			URL:     "https://da.wikipedia.org/wiki/Maskinl%C3%A6ring",
			Lang:    "da",
			Content: "Maskinlæring er en gren af kunstig intelligens, hvor algoritmer trænes til at finde mønstre i data.",
		},

	}

	for _, p := range seedData {
		if _, err := stmt.Exec(p.Title, p.URL, p.Lang, time.Now(), p.Content); err != nil {
			// Vi logger fejl, men fortsætter med næste række
			log.Printf("Error inserting seed data (%s): %v", p.Title, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
