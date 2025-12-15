import fs from "fs";

export function writePagesToJson(pages, path = "pages.json") {
    fs.writeFileSync(
        path,
        JSON.stringify(pages, null, 2),
        "utf8"
    );
}
