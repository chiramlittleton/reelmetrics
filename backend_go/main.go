package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/cors"
)

var ctx = context.Background()
var db *sql.DB
var redisClient *redis.Client

func init() {
	var err error

	db, err = sql.Open("pgx", "postgres://user:password@postgres:5432/reelmetrics_db?sslmode=disable")
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	log.Println("✅ Connected to PostgreSQL & Redis")
}

func getTheaters(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM theaters;")
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var theaters []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		theaters = append(theaters, map[string]interface{}{"id": id, "name": name})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(theaters)
}

func getMoviesByTheater(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	theaterID := vars["theater_id"]

	cacheKey := fmt.Sprintf("theater_%s_movies", theaterID)

	// Check Redis first
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"source": "cache", "data": %s}`, val)))
		return
	}

	// Query PostgreSQL if not cached
	query := `
		SELECT m.id, m.title, COALESCE(SUM(s.tickets_sold * s.ticket_price), 0) AS revenue
		FROM movies m
		LEFT JOIN sales s ON m.id = s.movie_id
		WHERE m.theater_id = $1
		GROUP BY m.id, m.title;
	`
	rows, err := db.Query(query, theaterID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var id int
		var title string
		var revenue float64
		rows.Scan(&id, &title, &revenue)
		movies = append(movies, map[string]interface{}{"id": id, "title": title, "ticket_sales": revenue})
	}

	// Store result in Redis
	jsonData, _ := json.Marshal(movies)
	redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// Get top theater by revenue
func getTopTheater(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	cacheKey := "top_theater_" + date

	// Check Redis first
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"source": "cache", "data": %s}`, val)))
		return
	}

	// Query PostgreSQL if not cached
	query := `
		SELECT t.name, COALESCE(SUM(s.tickets_sold * s.ticket_price), 0) AS total_revenue
		FROM sales s
		JOIN movies m ON s.movie_id = m.id
		JOIN theaters t ON m.theater_id = t.id
		WHERE s.sale_date = $1
		GROUP BY t.name
		ORDER BY total_revenue DESC
		LIMIT 1;
	`
	row := db.QueryRow(query, date)
	var theater string
	var revenue float64

	err = row.Scan(&theater, &revenue)
	if err != nil {
		http.Error(w, `{"error": "No data found"}`, http.StatusNotFound)
		return
	}

	// Store result in Redis
	result := map[string]interface{}{"theater": theater, "revenue": revenue}
	jsonData, _ := json.Marshal(result)
	redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	router := mux.NewRouter()

	// Define Routes
	router.HandleFunc("/theaters", getTheaters).Methods("GET")
	router.HandleFunc("/theaters/{theater_id}/movies", getMoviesByTheater).Methods("GET")
	router.HandleFunc("/top-theater/{date}", getTopTheater).Methods("GET")

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all (or restrict to ["http://localhost:3000"])
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
	}).Handler(router)

	log.Println("🚀 Go backend running on port 8002...")
	http.ListenAndServe(":8002", handler)
}
