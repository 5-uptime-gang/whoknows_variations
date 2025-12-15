import { normalizeQuery } from "./normalize.js";

export function clusterQueries(queries) {
    const topics = new Map();

    for (const q of queries) {
        const tokens = normalizeQuery(q);
        if (tokens.length < 2) continue;

        const candidates = new Set();

        // last 2 tokens
        candidates.add(tokens.slice(-2).join(" "));

        // last 3 tokens
        if (tokens.length >= 3) {
            candidates.add(tokens.slice(-3).join(" "));
        }

        // full tail noun phrase (up to 4 tokens)
        candidates.add(tokens.slice(-4).join(" "));

        for (const topic of candidates) {
            if (!topics.has(topic)) topics.set(topic, []);
            topics.get(topic).push(q);
        }
    }

    return topics;
}
