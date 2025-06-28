(() => {
  const page = document.body.dataset.page;
  if (page !== "register") {
    return;
  }

  const form = document.getElementById("registration-form") as HTMLFormElement;
  const errorMsg = document.getElementById("errorMsg") as HTMLSpanElement;

  errorMsg.classList.add("hidden");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(form);

    const name = formData.get("name");
    const email = formData.get("email");
    const password = formData.get("password");

    try {
      const res = await fetch("/register/create-user", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: name, email: email, password: password }),
      });

      if (!res.ok) {
        const errorText = await res.text();
        errorMsg.innerText = errorText;
        errorMsg.classList.remove("hidden");
      }
    } catch (err) {
      console.error("failed to register:", err);
    }
  });
})();
