import fs from "fs";

export function loadQueries(path) {
    return fs
        .readFileSync(path, "utf8")
        .split("\n")
        .map(q => q.trim())
        .filter(q => q.length > 3);
}
