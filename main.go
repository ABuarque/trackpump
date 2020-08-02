package main

import (
	"context"
	"log"
	"os"
	"trackpump/storage"

	"cloud.google.com/go/datastore"
	"github.com/labstack/echo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT environment variable")
	}
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("missing PROJECT_ID environment variable")
	}
	client, err := datastore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("falha ao criar cliente do Datastore, erro %q", err)
	}
	log.Println("connected to database")
	storageLogin := os.Getenv("STORAGE_LOGIN")
	if storageLogin == "" {
		log.Fatal("missing STORAGE_LOGIN environment variable")
	}
	storagePassword := os.Getenv("STORAGE_PASSWORD")
	if storagePassword == "" {
		log.Fatal("missing STORAGE_PASSWORD environment variable")
	}
	storageClient, err := storage.New(storageLogin, storagePassword)
	if err != nil {
		log.Fatalf("failed to create storage client, erro %q", err)
	}
	e := echo.New()
	userRegistry := NewRegistry(client, storageClient)
	usersControllers := userRegistry.NewAppController()
	e.POST("/api/v1/users", usersControllers.Create)
	e.POST("/api/v1/users/login", usersControllers.Login)
	e.GET("/api/v1/weekly_report", usersControllers.WeeklyReport)
	e.POST("/api/v1/measurements", usersControllers.RegisterMeasurement)
	e.GET("/", usersControllers.HomePage)
	e.GET("/sign_up", usersControllers.SignUp)
	e.POST("/process_signup", usersControllers.ProcessSignUp)
	e.GET("/login", usersControllers.LoginPage)
	e.POST("/process_login", usersControllers.ProcessLogin)
	e.GET("/admin", usersControllers.Admin)
	e.GET("/measurement", usersControllers.MeasurementPage)
	e.POST("/process_measurement", usersControllers.ProcessMeasurement)
	log.Println("server online at ", port)
	log.Fatal(e.Start(":" + port))
}
