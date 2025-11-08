document.getElementById("loginForm").addEventListener("submit", (e) => {
  e.preventDefault();
  const form = e.target;
  const formData = new FormData(form);
  const errorMessage = document.getElementById("errorMessage");

  let responseStatus; // Gem response status her

  fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      username: formData.get("username"),
      password: formData.get("password"),
    }),
  })
    .then((response) => {
      responseStatus = response.ok; // Gem response.ok
      return response.json();
    })
    .then((data) => {
      if (responseStatus) {
        // Brug den gemte status
        // Login successful
        window.location.href = "/";
      } else {
        // Show error message
        errorMessage.textContent = data.detail[0].msg || "Login failed";
        errorMessage.style.display = "block";
      }
    })
    .catch((error) => {
      errorMessage.textContent = "Error connecting to server";
      errorMessage.style.display = "block";
      console.error("Error:", error);
    });
});
