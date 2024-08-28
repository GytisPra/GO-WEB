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

	taskHandler := handlers.NewTaskHandler(taskService, sessionService)
	callbackHandler := middleware.NewCallbackHandler(userSerivce, accountService, sessionService)
	homeHandler := handlers.NewHomeHandler(sessionService)
	logoutHandler := handlers.NewLogoutHandler(sessionService)
	loginHandler := handlers.NewLoginHandler(sessionService)

	go sessionService.CleanupExpiredSessions()

	r.HandleFunc("/", homeHandler.HomeHandler)
	r.HandleFunc("/logout", logoutHandler.LogoutHandler)
	r.HandleFunc("/login", loginHandler.ShowLoginOptionsHandler)
	r.HandleFunc("/login/discord", loginHandler.LoginWithDiscordHandler)
	r.HandleFunc("/callback/discord", callbackHandler.DiscordCallbackHandler)
	r.HandleFunc("/task", taskHandler.ShowTaskFormHandler)
	r.HandleFunc("/task/all", taskHandler.ShowTasksHandler)
	r.HandleFunc("/task/create", taskHandler.CreateTaskHandler)

	log.Println("Server started on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
