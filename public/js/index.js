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
    const res = await fetch(url, {
      method: "GET",
    });

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

        wrapper.innerHTML = `
          <h2>
            <a class="search-result-title" href="${page.url}">
              ${page.title}
            </a>
          </h2>
          <p class="search-result-description">${page.content}</p>
          <p><strong>Language:</strong> ${page.language}</p>
          <p><strong>Last updated:</strong> ${new Date(
            page.last_updated
          ).toLocaleString()}</p>
        `;

        resultsContainer.appendChild(wrapper);
      });
    } else {
      resultsContainer.innerHTML = "<p>No results found.</p>";
    }
  } catch (err) {
    resultsContainer.innerHTML = `<p style="color:red;">Error: ${err.message}</p>`;
  }
}
