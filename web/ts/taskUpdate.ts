import { taskState } from "./state";

(() => {
  // Get modal elements
  const editModal = document.getElementById("editModal") as HTMLDivElement;
  const taskBody = document.getElementById("taskBody") as HTMLInputElement;
  const cancelBtn = document.getElementById("cancelEdit") as HTMLButtonElement;
  const saveBtn = document.getElementById("saveEdit") as HTMLButtonElement;
  const editErrorMsg = document.getElementById(
    "editErrorMsg"
  ) as HTMLSpanElement;

  // Attach click event to all edit buttons
  document.querySelectorAll(".edit-btn").forEach((btn) => {
    btn.addEventListener("click", (event) => {
      const target = event.currentTarget as HTMLElement;

      taskState.currentTaskBody = target.dataset.task || "";
      taskState.currentTaskId = target.dataset.id || "";

      taskBody.value = taskState.currentTaskBody;
      editModal.classList.remove("hidden");
    });
  });

  // Cancel button hides the modal & error
  cancelBtn.addEventListener("click", () => {
    editModal.classList.add("hidden");
    hideEditError();
  });

  // Save button submits the update
  saveBtn.addEventListener("click", async () => {
    const updatedBody = taskBody.value;

    if (!updatedBody.trim()) {
      showEditError("Task body cannot be empty");
      return;
    }

    if (!taskState.currentTaskId) {
      showEditError("Missing task ID");
      return;
    }

    try {
      const res = await fetch("/tasks/update", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          body: updatedBody,
          id: taskState.currentTaskId,
        }),
      });

      if (res.ok) {
        const editedTaskBody = document.getElementById(
          taskState.currentTaskId + "_taskBody"
        );
        if (editedTaskBody === null) {
          showEditError("Failed to find what task was just edited!");
          return;
        }

        editedTaskBody.innerText = updatedBody;

        const editButton = editedTaskBody
          .closest("li")
          ?.querySelector(".edit-btn") as HTMLElement;

        if (editButton) {
          editButton.dataset.task = updatedBody;
        }

        hideEditError();
      } else {
        showEditError("Failed to update task");
      }
    } catch (err) {
      console.error("Update failed:", err);
      showEditError("Error updating task");
    }

    editModal.classList.add("hidden");
  });

  function showEditError(msg: string) {
    editErrorMsg.classList.remove("hidden");
    editErrorMsg.innerText = msg;
  }

  function hideEditError() {
    if (!editErrorMsg.classList.contains("hidden")) {
      editErrorMsg.classList.add("hidden");
    }
    editErrorMsg.innerText = "";
  }
})();
