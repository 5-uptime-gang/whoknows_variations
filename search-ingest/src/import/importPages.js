import "dotenv/config";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { Pool } from "pg";

// Imports pages from a JSON file or STDIN into the database

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function getArg(flag, fallback = null) {
    const idx = process.argv.indexOf(flag);
    if (idx !== -1 && idx + 1 < process.argv.length) {
        return process.argv[idx + 1];
    }
    return fallback;
}

async function main() {
    const databaseUrl = getArg("--database-url", process.env.DATABASE_URL);
    const fileArg = getArg("--file");
    const filePath = fileArg
        ? path.resolve(fileArg)
        : path.join(__dirname, "../../pages.json");

    if (!databaseUrl) {
        throw new Error("DATABASE_URL is required (pass --database-url or set env)");
    }

    let raw;
    if (fileArg) {
        if (!fs.existsSync(filePath)) {
            throw new Error(`File not found: ${filePath}`);
        }
        raw = fs.readFileSync(filePath, "utf8");
        console.log(`Reading pages from file: ${filePath}`);
    } else {
        raw = fs.readFileSync(0, "utf8");
        console.log("Reading pages from STDIN");
    }

    if (!raw || raw.trim().length === 0) {
        throw new Error("Input is empty (no JSON data)");
    }

    const pages = JSON.parse(raw);
    if (!Array.isArray(pages) || pages.length === 0) {
        console.log("No pages to import");
        return;
    }

    const pool = new Pool({ connectionString: databaseUrl });
    const client = await pool.connect();

    try {
        await client.query("BEGIN");

        for (const page of pages) {
            await client.query(
                `
        INSERT INTO pages (title, url, language, last_updated, content)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (url) DO UPDATE
        SET title = EXCLUDED.title,
            language = EXCLUDED.language,
            last_updated = EXCLUDED.last_updated,
            content = EXCLUDED.content;
        `,
                [
                    page.title,
                    page.url,
                    page.language || "en",
                    page.last_updated || new Date().toISOString(),
                    page.content,
                ]
            );
        }

        await client.query("COMMIT");
        console.log(`Imported ${pages.length} pages`);
    } catch (err) {
        await client.query("ROLLBACK");
        throw err;
    } finally {
        client.release();
        await pool.end();
    }
}

main().catch((err) => {
    console.error("Import failed:", err.message);
    process.exit(1);
});
