package main

import (
	"log"
	"net/http"
	"os"
	"web-app/config"
	"web-app/internal/handlers"
	"web-app/internal/middleware"
	"web-app/internal/models"
	"web-app/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Handle(
	path string,
	r *mux.Router,
	handler http.HandlerFunc,
	middleware ...func(http.Handler) http.Handler,
) *mux.Route {
	// Chain middleware manually
	var h http.Handler = handler
	for _, m := range middleware {
		h = m(h)
	}
	return r.Handle(path, h)
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	log.Println("Starting app..")

	if err := config.ValidateEnv(); err != nil {
		log.Fatalf("Environment validation error: %v", err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Session{})
	db.AutoMigrate(&models.Task{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := mux.NewRouter()

	r.PathPrefix("/web/static/").Handler(http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static"))))

	taskService := services.NewTaskService(db)
	accountService := services.NewAcountService(db)
	userSerivce := services.NewUserService(db)
	sessionService := services.NewSessionService(db)

	authMiddleware := middleware.AuthMiddleware(sessionService)
	softAuthMiddleware := middleware.SoftAuthMiddleware(sessionService)

	taskHandler := handlers.NewTaskHandler(taskService, sessionService)
	callbackHandler := middleware.NewCallbackHandler(userSerivce, accountService, sessionService)
	homeHandler := handlers.NewHomeHandler(sessionService)
	logoutHandler := handlers.NewLogoutHandler(sessionService)
	loginHandler := handlers.NewLoginHandler(sessionService)

	go sessionService.CleanupExpiredSessions()

	// Don't use any auth middelware for these routes
	Handle("/login", r, loginHandler.ShowLoginOptionsHandler)
	Handle("/login/discord", r, loginHandler.LoginWithDiscordHandler)
	Handle("/callback/discord", r, callbackHandler.DiscordCallbackHandler)

	// Use a softAuthMiddleware for these routes so that we can check if the user is logged in (accesible to public)
	Handle("/", r, homeHandler.HomeHandler, softAuthMiddleware)

	// Protected routes only accesible when logged in
	Handle("/logout", r, logoutHandler.LogoutHandler, authMiddleware)
	Handle("/task", r, taskHandler.ShowTaskFormHandler, authMiddleware)
	Handle("/task/all", r, taskHandler.ShowTasksHandler, authMiddleware)
	Handle("/task/create", r, taskHandler.CreateTaskHandler, authMiddleware)

	log.Println("✅Server started. Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
