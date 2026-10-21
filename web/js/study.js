let currentCard

async function loadNextCard() {

  const params = new URLSearchParams(window.location.search)
  const deckID = params.get("deck_id")

  const response = await fetch(`/api/study/next?deck_id=${deckID}`)

  if (!response.ok) {
    console.log("failed to get next card")
    return
  }
  currentCard = await response.json()
  const frontContentDiv = document.getElementById("front-content")
  const backContentDiv = document.getElementById("back-content")

  frontContentDiv.innerHTML = `<p>${currentCard.front_content}</p>`
  backContentDiv.innerHTML = `<p>${currentCard.back_content}</p>`
}

async function evaluateCard(rating) {
  const params = new URLSearchParams(window.location.search)
  const deckID = params.get("deck_id")

  const response = await fetch("/api/study/evaluate", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      deck_id: deckID,
      card_id: currentCard.id,
      rating: rating,
    }),
  })

  if (!response.ok) {
    console.log("failed to evaluate card")
    return
  }

  await loadNextCard()
}

const ratingButtons = document.querySelectorAll(".rating-button")

for (const button of ratingButtons) {
  button.addEventListener("click", () => {
    const rating = Number(button.dataset.rating)
    evaluateCard(rating)
  })
}

loadNextCard()

