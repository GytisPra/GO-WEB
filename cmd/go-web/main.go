package main

import (
	"log"
	"net/http"
	"os"
	"web-app/internal/handlers"
	"web-app/internal/models"
	"web-app/internal/services"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		log.Fatal("DATABASE_URL is not set in the environment")
	}

	db, err := gorm.Open(postgres.Open(DATABASE_URL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database handle: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	http.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static"))))

	taskService := services.NewTaskService(db)
	taskHandler := handlers.NewTaskHandler(taskService)

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/task", taskHandler.ShowTasksHandler)
	http.HandleFunc("/task/create", taskHandler.CreateTaskHandler)

	log.Println("Server started on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
