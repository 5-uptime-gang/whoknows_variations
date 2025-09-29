let searchInput;

document.addEventListener("DOMContentLoaded", () => {
  searchInput = document.getElementById("search-input");

  // Focus the input field
  searchInput.focus();

  // Search when the user presses Enter
  searchInput.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
      // makeSearchRequest();
    }
  });
});
