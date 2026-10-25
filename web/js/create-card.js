const form = document.getElementById("card-form")

async function createCard(card) {
  const response = await fetch("/api/cards", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(card),
  })

  return response.ok
}

form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const params = new URLSearchParams(window.location.search)
  const deckID = params.get("deck_id")

  const selectedNotes = Array.from(
    document.querySelectorAll(".note-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedAccidentals = Array.from(
    document.querySelectorAll(".accidental-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedModes = Array.from(
    document.querySelectorAll(".mode-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedIntervals = Array.from(
    document.querySelectorAll(".interval-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedChords = Array.from(
    document.querySelectorAll(".chord-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedInversions = Array.from(
    document.querySelectorAll(".inversion-option:checked")
  ).map((checkbox) => checkbox.value)

  const selectedDirections = Array.from(
    document.querySelectorAll(".direction-option:checked")
  ).map((checkbox) => checkbox.value)

  const frontContent = document.getElementById("front-content").value
  const backContent = document.getElementById("back-content").value

  const replacements = [
    { tag: "{note}", values: selectedNotes },
    { tag: "{accidental}", values: selectedAccidentals },
    { tag: "{mode}", values: selectedModes },
    { tag: "{interval}", values: selectedIntervals },
    { tag: "{chord}", values: selectedChords },
    { tag: "{inversion}", values: selectedInversions },
    { tag: "{direction}", values: selectedDirections },
  ]

  const activeReplacements = replacements.filter((replacement) =>
    frontContent.includes(replacement.tag) ||
    backContent.includes(replacement.tag)
  )

  if (activeReplacements.length === 0) {
    const newCard = {
      front_content: frontContent,
      back_content: backContent,
      tempo: Number(document.getElementById("tempo").value),
      deck_id: deckID,
    }

    const success = await createCard(newCard)

    if (!success) {
      const error = document.getElementById("card-error")
      error.textContent = "Failed to create card"
      error.hidden = false
      return
    }

    alert("Card created successfully")
    window.location.href = `/edit-deck?deck_id=${deckID}`
  } else {
    for (const replacement of activeReplacements) {
      if (replacement.values.length === 0) {
        const error = document.getElementById("card-error")
        error.textContent = `Select at least one option for ${replacement.tag}`
        error.hidden = false
        return
      }
    }

    let generatedContent = [
      {
        front: frontContent,
        back: backContent,
      }
    ]

    for (const replacement of activeReplacements) {
      const nextGeneratedContent = []

      for (const content of generatedContent) {
        for (const value of replacement.values) {
          nextGeneratedContent.push({
            front: content.front.replaceAll(replacement.tag, value),
            back: content.back.replaceAll(replacement.tag, value),
          })
        }
      }

      generatedContent = nextGeneratedContent
    }

    const generatedCards = []

    for (const content of generatedContent) {
      generatedCards.push({
        front_content: content.front,
        back_content: content.back,
        tempo: Number(document.getElementById("tempo").value),
        deck_id: deckID,
      })
    }

    const cardCount = generatedCards.length
    const confirmed = confirm(`Create ${cardCount} cards?`)

    if (!confirmed) {
      return
    }

    for (const card of generatedCards) {
      const success = await createCard(card)

      if (!success) {
        const error = document.getElementById("card-error")
        error.textContent = "Failed to create generated cards"
        error.hidden = false
        return
      }
    }

    alert("Cards created successfully")
    window.location.href = `/edit-deck?deck_id=${deckID}`
  }
})

const advancedCheckbox = document.getElementById("advanced-options")
const advancedSection = document.getElementById("advanced-section")

advancedSection.style.display = "none"
advancedSection.style.alignItems = "flex-start"
advancedSection.style.gap = "24px"

advancedCheckbox.addEventListener("change", () => {
  if (advancedCheckbox.checked) {
    advancedSection.style.display = "flex"
  } else {
    advancedSection.style.display = "none"
  }
})

const notes = ["A", "B", "C", "D", "E", "F", "G"]

const accidentals = ["", "#", "b"]

const modes = [
  "Major",
  "Minor",
  "Harmonic Minor",
  "Melodic Minor",
  "Ionian",
  "Dorian",
  "Phrygian",
  "Lydian",
  "Mixolydian",
  "Aeolian",
  "Locrian",
]

const intervals = [
  "Unison",
  "Minor 2nd",
  "Major 2nd",
  "Augmented 2nd",
  "Minor 3rd",
  "Major 3rd",
  "Perfect 4th",
  "Augmented 4th",
  "Diminished 5th",
  "Perfect 5th",
  "Augmented 5th",
  "Minor 6th",
  "Major 6th",
  "Augmented 6th",
  "Diminished 7th",
  "Minor 7th",
  "Major 7th",
  "Octave",
  "Minor 9th",
  "Major 9th",
  "Augmented 9th",
  "Minor 10th",
  "Major 10th",
  "Perfect 11th",
  "Augmented 11th",
  "Perfect 12th",
  "Minor 13th",
  "Major 13th",
]

const inversions = [
  "Root",
  "1st",
  "2nd",
  "3rd",
]

const directions = [
  "Up",
  "Down",
]

const chords = [
  "Major",
  "Minor",
  "Diminished",
  "Augmented",
  "Major 7",
  "Minor 7",
  "Dominant 7",
  "Half-Diminished 7",
  "Diminished 7"
]

const noteOptions = document.getElementById("note-options")
const accidentalOptions = document.getElementById("accidental-options")
const modeOptions = document.getElementById("mode-options")
const intervalOptions = document.getElementById("interval-options")
const inversionOptions = document.getElementById("inversion-options")
const directionOptions = document.getElementById("direction-options")
const chordOptions = document.getElementById("chord-options")

for (const note of notes) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${note}" class="note-option">
      ${note}
    </label>
  `

  noteOptions.append(checkDiv)
}

for (const accidental of accidentals) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${accidental}" class="accidental-option">
      ${accidental}
    </label>
  `

  accidentalOptions.append(checkDiv)
}

for (const mode of modes) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${mode}" class="mode-option">
      ${mode}
    </label>
  `

  modeOptions.append(checkDiv)
}

for (const interval of intervals) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${interval}" class="interval-option">
      ${interval}
    </label>
  `

  intervalOptions.append(checkDiv)
}

for (const inversion of inversions) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${inversion}" class="inversion-option">
      ${inversion}
    </label>
  `

  inversionOptions.append(checkDiv)
}

for (const direction of directions) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${direction}" class="direction-option">
      ${direction}
    </label>
  `

  directionOptions.append(checkDiv)
}

for (const chord of chords) {
  const checkDiv = document.createElement("div")

  checkDiv.innerHTML = `
    <label>
      <input type="checkbox" value="${chord}" class="chord-option">
      ${chord}
    </label>
  `

  chordOptions.append(checkDiv)
}
