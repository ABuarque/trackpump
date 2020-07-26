package main

import (
	"fmt"
	"log"
	"os"
	"trackpump/db"
	"trackpump/storage"

	"github.com/labstack/echo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("missing PORT environment variable")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("missing DB_HOST environment variable")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("missing DB_NAME environment variable")
	}
	db, err := db.New(dbHost, dbName)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start db: %q", err))
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
	userRegistry := NewRegistry(db, storageClient)
	usersControllers := userRegistry.NewAppController()
	e.POST("/api/v1/users", usersControllers.Create)
	e.POST("/api/v1/users/login", usersControllers.Login)
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
