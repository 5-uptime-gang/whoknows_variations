import { removeStopwords } from "stopword";

export function normalizeQuery(query) {
    return removeStopwords(
        query
            .toLowerCase()
            .replace(/[^a-z0-9\s]/g, "")
            .split(/\s+/)
    );
}
