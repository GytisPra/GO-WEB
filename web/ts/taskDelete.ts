import { taskState } from "./state";

(() => {
  // Get modal elements
  const deleteModal = document.getElementById("deleteModal") as HTMLDivElement;
  const cancelDelete = document.getElementById(
    "cancelDelete"
  ) as HTMLButtonElement;
  const confirmDelete = document.getElementById(
    "confirmDelete"
  ) as HTMLButtonElement;

  // Attach click event to all delete buttons
  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", (event) => {
      const target = event.currentTarget as HTMLElement;

      taskState.currentTaskId = target.dataset.id || "";
      taskState.currentTaskBody = target.dataset.task || "";

      const taskBody = document.getElementById(
        "taskBody-delete"
      ) as HTMLSpanElement;

      if (taskState.currentTaskBody || taskState.currentTaskBody === "") {
        taskBody.innerText = taskState.currentTaskBody;
      } else {
        taskBody.innerText = "NOT_FOUND";
      }

      console.log("Delete button pressed!");

      deleteModal.classList.remove("hidden");
    });
  });

  // Cancel button hides the modal & error
  cancelDelete.addEventListener("click", () => {
    deleteModal.classList.add("hidden");
  });

  // Save button submits the update
  confirmDelete.addEventListener("click", async () => {
    if (!taskState.currentTaskId) {
      console.error("Missing task ID");
      return;
    }

    try {
      const res = await fetch("/tasks/delete", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: taskState.currentTaskId }),
      });

      if (res.ok) {
        const deletedTask = document.getElementById(
          taskState.currentTaskId + "_taskli"
        ) as HTMLLIElement;

        deletedTask.classList.add("hidden");
      }
    } catch (err) {
      console.error("failed to delete task:", err);
    }

    deleteModal.classList.add("hidden");
  });
})();
