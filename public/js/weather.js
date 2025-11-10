async function loadWeather() {
  const container = document.getElementById("weather-container");
  try {
    const res = await fetch("/api/weather");
    if (!res.ok) throw new Error("Failed to fetch weather");
    const result = await res.json();

    container.innerHTML = "";

    result.data.daily.forEach((day) => {
      const div = document.createElement("div");
      div.className = "weather-day";

      const dateSpan = document.createElement("span");
      dateSpan.className = "date";
      dateSpan.textContent = day.date;

      const tempSpan = document.createElement("span");
      tempSpan.className = "temps";
      tempSpan.textContent = `🌡 ${day.tMin}°C – ${day.tMax}°C`;

      const codeSpan = document.createElement("span");
      codeSpan.className = "code";
      codeSpan.textContent = `code ${day.code}`;

      div.appendChild(dateSpan);
      div.appendChild(tempSpan);
      div.appendChild(codeSpan);

      container.appendChild(div);
    });

    const footer = document.createElement("div");
    footer.className = "forecast-footer";
    footer.textContent = `Source: ${result.data.source} • Updated: ${new Date(
      result.data.updated
    ).toLocaleString()} (${result.data.timezone})`;
    container.appendChild(footer);
  } catch (err) {
    container.textContent = `Error loading weather data: ${err.message}`;
    container.style.color = "red";
  }
}

loadWeather();
