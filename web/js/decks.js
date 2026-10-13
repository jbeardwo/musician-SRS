
async function loadDecks() {
  const userID = localStorage.getItem("userID")
  const response = await fetch(`/api/decks?user_id=${userID}`)

  if (!response.ok) {
    console.log("failed to load decks")
    return
  }

  const decks = await response.json()
  const deckList = document.getElementById("deck-list")

  for (const deck of decks) {
    const deckDiv = document.createElement("div")

    deckDiv.innerHTML = `
      <h2>${deck.title}</h2>
      <p>${deck.description}</p>
      <button>Study</button>
      <a href="/edit-deck?deck_id=${deck.id}">Edit</a>
    `
    deckList.appendChild(deckDiv)
  }
}

loadDecks()
