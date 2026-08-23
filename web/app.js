const status = document.querySelector("#status");
const refresh = document.querySelector("#refresh");

async function loadStatus() {
  status.textContent = "Checking...";
  try {
    const response = await fetch("/healthz", { headers: { Accept: "application/json" } });
    const body = await response.json();
    status.textContent = response.ok ? `Healthy (${body.status})` : "Unavailable";
  } catch (error) {
    status.textContent = "Unavailable";
  }
}

refresh.addEventListener("click", loadStatus);
loadStatus();
