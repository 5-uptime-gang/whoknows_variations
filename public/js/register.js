document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector(".register-container");
  form.addEventListener("submit", (e) => {
    e.preventDefault();
    const formData = new FormData(form);
    const username = formData.get("username");
    const email = formData.get("email");
    const password = formData.get("password");
    const password2 = formData.get("password2");

    fetch("/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username,
        email,
        password,
        password2,
      }),
    })
      .then((response) => {
        if (response.status === 200) {
          return response.json();
        } else {
          return response.json().then((data) => {
            throw new Error(data.detail[0].msg || "Registration failed");
          });
        }
      })
      .then((data) => {
        // Registration successful, redirect to home page, user is logged in
        window.location.href = "/";
      })
      .catch((error) => {
        const errorMessageContainer = document.getElementById("errorMessage");
        errorMessageContainer.textContent = error.message;
        errorMessageContainer.style.display = "block";
      });
  });
});
