package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"web-app/internal/models"
	"web-app/internal/services"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) ShowTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var allTasks []models.Task
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/task-form.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}
	allTasks, err = h.taskService.GetAllTasks()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", allTasks)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/task", http.StatusSeeOther)
		return
	}

	var requestData struct {
		Body string `json:"body"`
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	requestData.Body = r.FormValue("task-body")

	_, err = h.taskService.CreateTask(requestData.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/task", http.StatusSeeOther)
}
