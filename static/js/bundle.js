"use strict";
(() => {
  // web/ts/taskDelete.ts
  var deleteModal = document.getElementById("deleteModal");
  var cancelDelete = document.getElementById(
    "cancelDelete"
  );
  var confirmDelete = document.getElementById(
    "confirmDelete"
  );
  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", (event) => {
      const target = event.currentTarget;
      currentTaskId = target.dataset.id || "";
      deleteModal.classList.remove("hidden");
    });
  });
  cancelDelete.addEventListener("click", () => {
    deleteModal.classList.add("hidden");
  });
  confirmDelete.addEventListener("click", async () => {
    if (!currentTaskId) {
      console.error("Missing task ID");
      return;
    }
    try {
      const res = await fetch("/tasks/delete", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id: currentTaskId })
      });
      if (res.ok) {
        const deletedTask = document.getElementById(
          currentTaskId + "_taskBody"
        );
        deletedTask.classList.add("hidden");
      }
    } catch (err) {
      console.error("failed to delete task:", err);
    }
    editModal.classList.add("hidden");
  });

  // web/ts/taskUpdate.ts
  var editModal2 = document.getElementById("editModal");
  var taskBody = document.getElementById("taskBody");
  var cancelBtn = document.getElementById("cancelEdit");
  var saveBtn = document.getElementById("saveEdit");
  var editErrorMsg = document.getElementById("editErrorMsg");
  document.querySelectorAll(".edit-btn").forEach((btn) => {
    btn.addEventListener("click", (event) => {
      const target = event.currentTarget;
      currentTaskBody = target.dataset.task || "";
      currentTaskId = target.dataset.id || "";
      taskBody.value = currentTaskBody;
      editModal2.classList.remove("hidden");
    });
  });
  cancelBtn.addEventListener("click", () => {
    editModal2.classList.add("hidden");
    hideEditError();
  });
  saveBtn.addEventListener("click", async () => {
    const updatedBody = taskBody.value;
    if (!updatedBody.trim()) {
      showEditError("Task body cannot be empty");
      return;
    }
    if (!currentTaskId) {
      showEditError("Missing task ID");
      return;
    }
    try {
      const res = await fetch("/tasks/update", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ body: updatedBody, id: currentTaskId })
      });
      if (res.ok) {
        const editedTaskBody = document.getElementById(
          currentTaskId + "_taskBody"
        );
        if (editedTaskBody === null) {
          showEditError("Failed to find what task was just edited!");
          return;
        }
        editedTaskBody.innerText = updatedBody;
        const editButton = editedTaskBody.closest("li")?.querySelector(".edit-btn");
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
    editModal2.classList.add("hidden");
  });
  function showEditError(msg) {
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
