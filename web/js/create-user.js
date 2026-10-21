const form = document.getElementById("user-form")


form.addEventListener("submit", async (event) => {
  event.preventDefault()

  const email = document.getElementById("email").value

  const response = await fetch("/api/users", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      email: email,
    }),
  })

  if (!response.ok) {
    const error = document.getElementById("user-error")
    error.textContent = "Failed to create user"
    error.hidden = false
    return
  }

  alert("User created successfully")
  window.location.href = "/"
})
