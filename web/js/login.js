
const form = document.getElementById("login-form")
const emailInput = document.getElementById("email")

form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const email = emailInput.value

  const response = await fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: email,
    }),
  })

  if (!response.ok) {
    const error = document.getElementById("login-error")
    error.textContent = "User not found"
    error.hidden = false
  } else {
    const user = await response.json()
    localStorage.setItem("userID", user.id)
    window.location.href = "/decks"
  }
})
