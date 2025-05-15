package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arjunsaxaena/MakerbleAssignment/patient_service/controller"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/database"
	"github.com/arjunsaxaena/MakerbleAssignment/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	PORT = 8002
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}

	database.Connect()
	defer database.Close()

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	patientController := controller.NewPatientController()
	
	patientRoutes := router.Group("/api/patients")
	patientRoutes.Use(middleware.AuthMiddleware())
	{
		patientRoutes.POST("", middleware.RoleAuthorization("receptionist"), patientController.Create)
		patientRoutes.DELETE("/:id", middleware.RoleAuthorization("receptionist"), patientController.Delete)
		
		patientRoutes.PATCH("", middleware.RoleAuthorization("doctor"), patientController.Update)
		
		patientRoutes.GET("", middleware.RoleAuthorization("doctor", "receptionist"), patientController.Get)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: router,
	}

	go func() {
		log.Printf("Patient service listening on port %d\n", PORT)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
} 