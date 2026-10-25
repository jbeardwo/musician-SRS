const form = document.getElementById("deck-form")
form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const userID = localStorage.getItem("userID")

  const newDeck = {
    title: document.getElementById("deck-title").value,
    description: document.getElementById("deck-description").value,
    user_id: userID,
  }

  const response = await fetch("/api/decks", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(newDeck),
  })

  if (!response.ok) {
    const error = document.getElementById("deck-error")
    error.textContent = "Failed to create deck"
    error.hidden = false
    return
  }

  alert("Deck created successfully")
  window.location.href = "/decks"
})
