async function checkSession() {
  try {
    const res = await fetch("/api/session", { credentials: "include" });
    const data = await res.json();

    if (data.logged_in) {
      document.getElementById("nav-logout").style.display = "inline-block";
      document.getElementById("nav-login").style.display = "none";
      document.getElementById("nav-register").style.display = "none";
    } else {
      document.getElementById("nav-logout").style.display = "none";
      document.getElementById("nav-login").style.display = "inline-block";
      document.getElementById("nav-register").style.display = "inline-block";
    }
  } catch (err) {
    console.error("session check failed", err);
  }

  // logout button click
  document.getElementById("nav-logout").addEventListener("click", async (e) => {
    e.preventDefault();
    await fetch("/api/logout", { method: "GET", credentials: "include" });
    location.reload();
  });
}

checkSession();
