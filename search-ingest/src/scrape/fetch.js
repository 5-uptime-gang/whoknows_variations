import axios from "axios";

export async function fetchHtml(url) {
    const res = await axios.get(url, {
        timeout: 15000,
        headers: {
            "User-Agent": "search-ingest/1.0"
        }
    });
    return res.data;
}
