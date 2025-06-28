import { taskState } from "./state";

(() => {
  const page = document.body.dataset.page;
  if (page !== "tasks") {
    return;
  }

  // Get modal elements
  const deleteModal = document.getElementById("deleteModal") as HTMLDivElement;
  const cancelDelete = document.getElementById(
    "cancelDelete",
  ) as HTMLButtonElement;
  const confirmDelete = document.getElementById(
    "confirmDelete",
  ) as HTMLButtonElement;
  const allTasks = document.getElementById("allTasks") as HTMLUListElement;

  // Attach click event to all delete buttons
  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", (event) => {
      const target = event.currentTarget as HTMLElement;

      taskState.currentTaskId = target.dataset.id || "";
      taskState.currentTaskBody = target.dataset.task || "";

      deleteModal.classList.add("flex");
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
          taskState.currentTaskId + "_taskli",
        ) as HTMLLIElement;

        deletedTask.classList.add("hidden");

        const remainingTasks =
          allTasks.querySelectorAll("li:not(.hidden)").length;

        if (remainingTasks == 0) {
          const allTasksDeleteP = document.getElementById(
            "allTasksDeleted",
          ) as HTMLParagraphElement;
          allTasksDeleteP.classList.add("flex");
          allTasksDeleteP.classList.remove("hidden");
        }
      }
    } catch (err) {
      console.error("failed to delete task:", err);
    }

    deleteModal.classList.remove("flex");
    deleteModal.classList.add("hidden");
  });
})();
