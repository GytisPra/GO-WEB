package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"web-app/internal/models"
	"web-app/internal/services"
)

type TaskHandler struct {
	taskService    *services.TaskService
	sessionService *services.SessionService
}

func NewTaskHandler(taskService *services.TaskService, sessionService *services.SessionService) *TaskHandler {
	return &TaskHandler{taskService: taskService, sessionService: sessionService}
}

func (h *TaskHandler) ShowTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isLoggedIn, user, err := h.sessionService.IsUserLoggedIn(cookie.Value)
	if err != nil || !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/task-view.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	allTasks, err := h.taskService.GetUserTasks(user.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}

	type ResponseData struct {
		AllTasks   []models.Task
		User       models.User
		IsLoggedIn bool
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", ResponseData{AllTasks: allTasks, User: *user, IsLoggedIn: isLoggedIn})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) ShowTaskFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	isLoggedIn, user, err := h.sessionService.IsUserLoggedIn(cookie.Value)
	if err != nil || !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/task-form.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template file: %v", err), http.StatusInternalServerError)
		return
	}

	type ResponseData struct {
		User       models.User
		IsLoggedIn bool
	}

	err = tmpl.ExecuteTemplate(w, "layout.html", ResponseData{User: *user, IsLoggedIn: isLoggedIn})
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
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	isLoggedIn, user, err := h.sessionService.IsUserLoggedIn(cookie.Value)
	if err != nil || !isLoggedIn {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	requestData.UserID = user.ID

	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	requestData.Body = r.FormValue("task-body")

	_, err = h.taskService.CreateTask(requestData.Body, requestData.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/task", http.StatusSeeOther)
}
