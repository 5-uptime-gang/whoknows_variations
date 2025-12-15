import "dotenv/config";
import pLimit from "p-limit";

import { loadQueries } from "./ingest/loadQueries.js";
import { clusterQueries } from "./ingest/cluster.js";
import { urlsForTopic } from "./ingest/topics.js";
import { scrape } from "./scrape/scrape.js";
import { writePagesToJson } from "./output/writeJson.js";

const limit = pLimit(4);
const MAX_URLS_PER_TOPIC = 5;

async function run() {
    const queries = loadQueries("queries.txt");
    const topics = clusterQueries(queries);

    console.log(`Found ${topics.size} topics`);

    const pages = [];
    const seenUrls = new Set();

    for (const [topic, qs] of topics.entries()) {
        if (qs.length < 3) continue; // ignore weak topics

        const urls = urlsForTopic(topic).slice(0, MAX_URLS_PER_TOPIC);

        for (const url of urls) {
            if (seenUrls.has(url)) continue;
            seenUrls.add(url);

            await limit(async () => {
                try {
                    const page = await scrape(url);

                    pages.push({
                        title: page.title,
                        url: page.url,
                        language: page.language,
                        last_updated: new Date().toISOString(),
                        content: page.content
                    });

                    console.log(`Collected: ${topic} → ${url}`);
                } catch {
                    console.log(`Skipped: ${url}`);
                }
            });
        }
    }

    writePagesToJson(pages, "pages.json");
    console.log(`\nSaved ${pages.length} pages to pages.json`);

    process.exit(0);
}

run();
