package main

import (
	"GoMin/config"
	"GoMin/handlers"
	"GoMin/miniohelper"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.LoadConfig()

	minioHelper, err := miniohelper.NewMinioHelper(
		config.AppConfig.Minio.Endpoint,
		config.AppConfig.Minio.AccessKey,
		config.AppConfig.Minio.SecretKey,
		config.AppConfig.Minio.UseSSL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	apiHandler := handlers.NewAPIHandler(minioHelper)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/v1/upload/:bucketName/file", apiHandler.UploadFile)
	e.DELETE("/v1/delete/:bucketName/file", apiHandler.RemoveFile)
	e.GET("/v1/retrieve/:bucketName/file", apiHandler.GetFile)

	port := config.AppConfig.Server.Port
	log.Printf("Starting server on port %d...\n", port)
	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
