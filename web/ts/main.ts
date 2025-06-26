// Get modal elements
const editModal = document.getElementById("editModal") as HTMLDivElement;
const editInput = document.getElementById("editInput") as HTMLInputElement;
const cancelBtn = document.getElementById("cancelEdit") as HTMLButtonElement;
const saveBtn = document.getElementById("saveEdit") as HTMLButtonElement;

// State: store currently editing task
let currentTaskBody: string | null = null;

// Attach click event to all edit buttons
document.querySelectorAll(".edit-btn").forEach((btn) => {
  btn.addEventListener("click", (event) => {
    const target = event.currentTarget as HTMLElement;
    currentTaskBody = target.dataset.task || "";
    editInput.value = currentTaskBody;
    editModal.classList.remove("hidden");
  });
});

// Cancel button hides the modal
cancelBtn.addEventListener("click", () => {
  editModal.classList.add("hidden");
});

// Save button submits the update
saveBtn.addEventListener("click", async () => {
  const updatedBody = editInput.value;

  if (!updatedBody.trim()) {
    alert("Task cannot be empty.");
    return;
  }

  // Example: Send to server (adjust this as needed)
  try {
    const res = await fetch("/tasks/update", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ oldBody: currentTaskBody, newBody: updatedBody }),
    });

    if (res.ok) {
      location.reload(); // Or update the DOM directly
    } else {
      alert("Failed to update task");
    }
  } catch (err) {
    console.error("Update failed:", err);
    alert("Error updating task");
  }

  editModal.classList.add("hidden");
});
