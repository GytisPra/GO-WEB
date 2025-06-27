package task

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"web-app/internal/middleware"
	"web-app/internal/models"
	"web-app/internal/services"
	"web-app/pkg/utils"
)

type TaskHandler struct {
	taskService    *services.TaskService
	sessionService *services.SessionService
}

func NewTaskHandler(taskService *services.TaskService, sessionService *services.SessionService) *TaskHandler {
	return &TaskHandler{taskService: taskService, sessionService: sessionService}
}

func (h *TaskHandler) ShowAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.FromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	allTasks, err := h.taskService.GetUserTasks(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}

	pageData := struct {
		AllTasks []models.Task
	}{AllTasks: allTasks}

	tmplData := utils.TemplateData{
		BuildTime:  utils.BuildTime,
		IsLoggedIn: true,
		Data:       pageData,
	}

	utils.RenderTemplate(w, "task-view.html", tmplData)
}

func (h *TaskHandler) ShowCreateTaskForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, ok := middleware.FromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pageData := struct{}{}

	tmplData := utils.TemplateData{
		BuildTime:  utils.BuildTime,
		IsLoggedIn: true,
		Data:       pageData,
	}

	utils.RenderTemplate(w, "task-form.html", tmplData)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/task", http.StatusSeeOther)
		return
	}

	var requestData struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	userID, ok := middleware.FromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	requestData.UserID = userID

	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	requestData.Body = r.FormValue("task-body")

	for i := 0; i < 100; i++ {
		_, err := h.taskService.CreateTask(requestData.Body+strconv.Itoa(i), requestData.UserID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/tasks/all", http.StatusSeeOther)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.FromContext(r.Context())
	if !ok {
		utils.HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var requestData struct {
		Body string `json:"body"`
		ID   string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.HandleError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTaskById(requestData.ID)
	if err != nil {
		utils.HandleError(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.UserID != userID {
		utils.HandleError(w, "Forbidden: You can only edit your own tasks", http.StatusForbidden)
		return
	}

	err = h.taskService.UpdateTask(userID, requestData.ID, requestData.Body)
	if err != nil {
		utils.HandleError(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.FromContext(r.Context())
	if !ok {
		utils.HandleError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var requestData struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.HandleError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTaskById(requestData.ID)
	if err != nil {
		utils.HandleError(w, "Task not found", http.StatusNotFound)
		return
	}

	if task.UserID != userID {
		utils.HandleError(w, "Forbidden: You can only delete your own tasks", http.StatusForbidden)
		return
	}

	if err := h.taskService.DeleteTask(task.ID); err != nil {
		utils.HandleError(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
