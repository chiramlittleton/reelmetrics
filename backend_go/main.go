package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"backend/db"
	"backend/handlers"
	"backend/redis"
)

func main() {
	// Initialize database & Redis
	db.ConnectDB()
	redis.ConnectRedis()

	// Load initial sales data into Redis
	redis.InitializeRedisCache()

	// Set up routes
	router := mux.NewRouter()
	router.HandleFunc("/theaters", handlers.GetTheaters).Methods("GET")
	router.HandleFunc("/theaters/{theater_id}/movies", handlers.GetMoviesByTheater).Methods("GET")
	router.HandleFunc("/top-theater/{date}", handlers.GetTopTheater).Methods("GET")

	// Start server
	log.Println("🚀 Go backend running on port 8002...")
	http.ListenAndServe(":8002", cors.AllowAll().Handler(router))
}
