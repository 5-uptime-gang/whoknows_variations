import { sample } from "../utils/random.js";

const MAX_URLS_PER_TOPIC = 5;

export function urlsForTopic(topic) {
    const slug = topic.replace(/\s+/g, "_");
    const encoded = encodeURIComponent(topic);

    const candidates = [
        // Encyclopedic / reference
        `https://en.wikipedia.org/wiki/${slug}`,
        `https://www.britannica.com/search?query=${encoded}`,

        // Programming & CS tutorials
        `https://www.geeksforgeeks.org/${encoded}/`,
        `https://www.javatpoint.com/search?query=${encoded}`,
        `https://www.tutorialspoint.com/search/${encoded}`,
        `https://www.programiz.com/search?q=${encoded}`,

        // Documentation hubs
        `https://www.ibm.com/topics/${slug}`,
        `https://learn.microsoft.com/en-us/search/?terms=${encoded}`,
        `https://developer.mozilla.org/en-US/search?q=${encoded}`,

        // Q&A / explanations
        `https://stackoverflow.com/search?q=${encoded}`,
        `https://stackprinter.appspot.com/v/1/search?q=${encoded}`,

        // Academic / learning
        `https://www.geeksforgeeks.org/?s=${encoded}`,
        `https://www.cs.cmu.edu/search/?q=${encoded}`,

        // Blog aggregators
        `https://medium.com/search?q=${encoded}`,
        `https://dev.to/search?q=${encoded}`,

        // Vendor-neutral
        `https://www.redhat.com/en/topics/${slug}`,
        `https://www.cloudflare.com/learning/${slug}/`
    ];

    // Randomly sample up to MAX_URLS_PER_TOPIC
    return sample(candidates, MAX_URLS_PER_TOPIC);
}
