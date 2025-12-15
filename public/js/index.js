document.addEventListener("DOMContentLoaded", async () => {
  const input = document.getElementById("search-input");
  const button = document.getElementById("search-button");

  input?.focus();

  if (button && input) {
    button.addEventListener("click", () => {
      doSearch(input.value);
    });

    input.addEventListener("keypress", (event) => {
      if (event.key === "Enter") {
        doSearch(input.value);
      }
    });
  }
});

async function doSearch(query, language = null) {
  const resultsContainer = document.getElementById("results");

  try {
    // Construct query string
    let url = `/api/search?q=${encodeURIComponent(query)}`;
    if (language) {
      url += `&language=${encodeURIComponent(language)}`;
    }

    // GET request
    const res = await fetch(url, { method: "GET" });

    if (!res.ok) {
      throw new Error(`Search failed with status ${res.status}`);
    }

    const data = await res.json();

    // Clear any old results
    resultsContainer.innerHTML = "";

    // Fill results
    if (data.data && Array.isArray(data.data) && data.data.length > 0) {
      data.data.forEach((page) => {
        const wrapper = document.createElement("div");

        // Title
        const h2 = document.createElement("h2");
        const link = document.createElement("a");
        link.className = "search-result-title";
        link.setAttribute("href", page.url);
        link.textContent = page.title;
        h2.appendChild(link);

        // Content snippet
        const desc = document.createElement("p");
        desc.className = "search-result-description";
        const snippet = page.snippet || page.content || "";
        desc.innerHTML = snippet;

        // Language
        const lang = document.createElement("p");
        lang.textContent = `Language: ${page.language}`;

        // Last updated
        const updated = document.createElement("p");
        updated.textContent = `Last updated: ${new Date(
          page.last_updated
        ).toLocaleString()}`;

        wrapper.appendChild(h2);
        wrapper.appendChild(desc);
        wrapper.appendChild(lang);
        wrapper.appendChild(updated);

        resultsContainer.appendChild(wrapper);
      });
    } else {
      const noResults = document.createElement("p");
      noResults.textContent = "No results found.";
      resultsContainer.appendChild(noResults);
    }
  } catch (err) {
    resultsContainer.innerHTML = "";
    const errorP = document.createElement("p");
    errorP.style.color = "red";
    errorP.textContent = `Error: ${err.message}`;
    resultsContainer.appendChild(errorP);
  }
}
