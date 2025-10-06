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

// isLive is a flag that indicates if the server is live.
var isLive = true

// main is the entry point of the application.
// It loads the .env file, checks if the "health" argument is passed, and runs the server.
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

// healthCheck performs a health check on the server.
// It sends a GET request to the /health/liveness endpoint and checks if the status is OK.
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

// runServer starts the Gin server.
// It sets up the routes and starts the server on the specified port.
// It also handles graceful shutdown of the server.
func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	r := gin.Default()

	// @Summary Root endpoint
	// @Description Returns a simple "Hello from Gin!" message.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
		})
	})

	// @Summary Liveness probe
	// @Description Returns the liveness status of the server.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Failure 503 {object} gin.H
	// @Router /health/liveness [get]
	r.GET("/health/liveness", func(c *gin.Context) {
		if isLive {
			c.JSON(http.StatusOK, gin.H{"status": "UP"})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		}
	})

	// @Summary Readiness probe
	// @Description Returns the readiness status of the server.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Router /health/readiness [get]
	r.GET("/health/readiness", func(c *gin.Context) {
		// In a real application, you would check database connections, external services, etc.
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	// @Summary Toggle liveness
	// @Description Toggles the liveness status of the server.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Router /health/liveness/toggle [post]
	r.POST("/health/liveness/toggle", func(c *gin.Context) {
		isLive = !isLive
		c.JSON(http.StatusOK, gin.H{"is_live": isLive})
	})

	// @Summary Get environment variable
	// @Description Returns the value of the CURR_ENV environment variable.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Router /env [get]
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

	// @Summary Shutdown server
	// @Description Initiates a graceful shutdown of the server.
	// @Produce json
	// @Success 200 {object} gin.H
	// @Router /shutdown [post]
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
