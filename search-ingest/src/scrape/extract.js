import { JSDOM } from "jsdom";
import { Readability } from "@mozilla/readability";

export function extractContent(html, url) {
    const dom = new JSDOM(html, { url });
    const reader = new Readability(dom.window.document);
    const article = reader.parse();

    if (!article?.textContent) {
        throw new Error("No readable content");
    }

    return {
        title: article.title,
        content: article.textContent.trim()
    };
}
