let currentCard

async function loadCard() {

  const params = new URLSearchParams(window.location.search)
  const cardID = params.get("card_id")
  const deckID = params.get("deck_id")

  document.getElementById("back").href = `/edit-deck?deck_id=${deckID}`

  const response = await fetch(`/api/cards/${cardID}`)


  if (!response.ok) {
    console.log("failed to load card")
    return
  }

  currentCard = await response.json()

  document.getElementById("front-content").value = currentCard.front_content
  document.getElementById("back-content").value = currentCard.back_content
}

const form = document.getElementById("card-form")

form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const updatedCard = {
    id: currentCard.id,
    front_content: document.getElementById("front-content").value,
    back_content: document.getElementById("back-content").value,
    interval: currentCard.interval,
    target: currentCard.target,
    ease_factor: currentCard.ease_factor,
    repetitions_count: currentCard.repetitions_count,
    last_reviewed_at: currentCard.last_reviewed_at,
    last_reviewed_num: currentCard.last_reviewed_num,
    tempo: currentCard.tempo,
    mastered_tempo: currentCard.mastered_tempo,
    mastery_stage: currentCard.mastery_stage,
    mastery_stage_time_started: currentCard.mastery_stage_time_started,
    challenge_again_count: currentCard.challenge_again_count
  }

  const response = await fetch("/api/cards", {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(updatedCard),
  })

  const error = document.getElementById("card-error")

  if (!response.ok) {
    error.textContent = "Card not found"
  } else {
    currentCard = await response.json()
    error.textContent = "Saved"
    error.hidden = false
  }
})

loadCard()
