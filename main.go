package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "health" {
		healthCheck()
		return
	}

	runServer()
}

func healthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	resp, err := http.Get("http://localhost:" + port + "/health/liveness")
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Health check failed with status %s", resp.Status)
	}

	log.Println("Health check successful")
}

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
		})
	})

	r.GET("/health/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	r.GET("/health/readiness", func(c *gin.Context) {
		// In a real application, you would check database connections, external services, etc.
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	r.GET("/env", func(c *gin.Context) {
		currEnv := os.Getenv("CURR_ENV")
		c.JSON(http.StatusOK, gin.H{"curr_env": currEnv})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	r.POST("/shutdown", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is shutting down gracefully..."})
		go func() {
			time.Sleep(100 * time.Millisecond)
			quit <- syscall.SIGTERM
		}()
	})

	go func() {
		log.Printf("Server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
