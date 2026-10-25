async function loadCards() {

  const params = new URLSearchParams(window.location.search)
  const deckID = params.get("deck_id")

  document.getElementById("back").href = "/decks"

  const response = await fetch(`/api/cards?deck_id=${deckID}`)

  if (!response.ok) {
    console.log("failed to load cards")
    return
  }

  const cards = await response.json()
  const cardList = document.getElementById("card-list")
  for (const card of cards) {
    const cardDiv = document.createElement("div")

    cardDiv.innerHTML = `
      <h2>${card.front_content}</h2>
      <p>${card.back_content}</p>
      <a href="/edit-card?card_id=${card.id}&deck_id=${deckID}">Edit</a>
    `
    cardList.appendChild(cardDiv)
  }
}

let currentDeck

async function loadDeck() {
  const params = new URLSearchParams(window.location.search)
  const deckID = params.get("deck_id")

  document.getElementById("create-card").href = `/create-card?deck_id=${deckID}`

  const response = await fetch(`/api/decks/${deckID}`)

  if (!response.ok) {
    console.log("failed to load deck")
    return
  }

  currentDeck = await response.json()

  document.getElementById("deck-title").value = currentDeck.title
  document.getElementById("deck-description").value = currentDeck.description
}

const form = document.getElementById("deck-form")

form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const updatedDeck = {
    id: currentDeck.id,
    title: document.getElementById("deck-title").value,
    description: document.getElementById("deck-description").value,
    total_reviews: currentDeck.total_reviews,
  }

  const response = await fetch("/api/decks", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updatedDeck),
  })

  const error = document.getElementById("deck-error")
  if (!response.ok) {
    error.textContent = "Deck not found"
    error.hidden = false
  } else {
    const currentDeck = await response.json()
    error.textContent = "Saved"
    error.hidden = false
  }
})

loadDeck()
loadCards()
