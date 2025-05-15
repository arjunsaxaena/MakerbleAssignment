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

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/database"
	"github.com/arjunsaxaena/MakerbleAssignment/portal_service/controller"
	"github.com/gin-gonic/gin"
)

const (
	PORT = 8001
)

func main() {
	database.Connect()
	defer database.Close()

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	userController := controller.NewUserController()
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("", userController.Create)
		userRoutes.GET("", userController.Get)
		userRoutes.PATCH("", userController.Update)
		userRoutes.DELETE("/:id", userController.Delete)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", PORT),
		Handler: router,
	}

	go func() {
		log.Printf("Portal service listening on port %d\n", PORT)
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
