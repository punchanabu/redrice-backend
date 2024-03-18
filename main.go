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

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	config "github.com/punchanabu/redrice-backend-go/config"
	routers "github.com/punchanabu/redrice-backend-go/routers"
	"github.com/punchanabu/redrice-backend-go/routers/api"
	v1 "github.com/punchanabu/redrice-backend-go/routers/api/v1"
)

func main() {

	// Load the env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup database connection
	db := config.SetupDBConnection()
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	// Initialize necessary handlers
	v1.InitializedUserHandler(db)
	v1.InitializedRestaurantHandler(db)
	api.InitializedAuthHandler(db)
	v1.InitializedReservationHandler(db)

	// Initialize router
	r := routers.UseRouter()
	r.Use(cors.New(config.CORSConfig()))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	fmt.Println("Gracefully shutting down server..., press Ctrl+C again to force shutdown")
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

}
