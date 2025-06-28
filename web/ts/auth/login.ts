(() => {
  const page = document.body.dataset.page;
  if (page !== "login") {
    return;
  }

  const form = document.getElementById("login-form") as HTMLFormElement;
  const errorMsg = document.getElementById("errorMsg") as HTMLSpanElement;

  errorMsg.classList.add("hidden");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    errorMsg.classList.add("hidden");

    const formData = new FormData(form);

    const email = formData.get("email");
    const password = formData.get("password");

    try {
      const res = await fetch("/login/local", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: email, password: password }),
      });

      if (res.redirected) {
        window.location.href = res.url;
        return;
      }

      if (!res.ok) {
        const errorText = await res.text();
        errorMsg.innerText = errorText;
        errorMsg.classList.remove("hidden");
      }
    } catch (err) {
      console.error("failed to login:", err);
    }
  });
})();
