document.addEventListener("DOMContentLoaded", () => {
  const checkDatesBtn = document.getElementById("check-dates-btn");
  const resultsContainer = document.getElementById("results-container");
  const loader = document.getElementById("loader");

  checkDatesBtn.addEventListener("click", async () => {
    // Show loader and clear previous results
    loader.classList.remove("hidden");
    resultsContainer.innerHTML = "";
    checkDatesBtn.disabled = true;

    try {
      // Fetch data from our Go backend API
      const response = await fetch("/api/available-dates");
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();

      // Sort city names alphabetically for consistent order
      const sortedCities = Object.keys(data).sort(
        (a, b) => data[b].dates.length - data[a].dates.length,
      );

      // Loop through the sorted city names and create cards
      for (const cityName of sortedCities) {
        const result = data[cityName];
        const card = document.createElement("div");
        card.classList.add("city-card");

        let content = `<h2>${result.centerName}</h2>`;

        if (result.error) {
          card.classList.add("unavailable");
          content += `<p class="error">Error: ${result.error}</p>`;
        } else if (result.dates && result.dates.length > 0) {
          card.classList.add("available");
          content += `<p>Found ${result.dates.length} available date(s):</p>`;
          content += "<ul>";
          result.dates.forEach((date) => {
            const [day, month, year] = date.bookingDate.split("-");
            const parsedDate = new Date(`${year}-${month}-${day}T00:00:00`);

            const formattedDate = parsedDate.toLocaleString("en-US", {
              year: "numeric",
              month: "long",
              day: "numeric",
              hour: "2-digit",
              minute: "2-digit",
            });
            content += `<li>${formattedDate}</li>`;
          });
          content += "</ul>";
        } else {
          card.classList.add("unavailable");
          content += "<p>No available dates found.</p>";
        }

        card.innerHTML = content;
        resultsContainer.appendChild(card);
      }
    } catch (error) {
      resultsContainer.innerHTML = `<p style="color: red; grid-column: 1 / -1;">Failed to fetch data: ${error.message}</p>`;
      console.error("Fetch error:", error);
    } finally {
      // Hide loader and re-enable button
      loader.classList.add("hidden");
      checkDatesBtn.disabled = false;
    }
  });
});
