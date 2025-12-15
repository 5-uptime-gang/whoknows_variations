import detect from "langdetect";
import { fetchHtml } from "./fetch.js";
import { extractContent } from "./extract.js";

export async function scrape(url) {
    const html = await fetchHtml(url);
    const { title, content } = extractContent(html, url);

    const detected = detect.detectOne(content);
    const language = detected === "da" ? "da" : "en";

    return { title, url, content, language };
}
