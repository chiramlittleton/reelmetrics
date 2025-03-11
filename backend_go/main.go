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
	_ "github.com/jackc/pgx/v5/stdlib"
)

var ctx = context.Background()
var db *sql.DB
var redisClient *redis.Client

func init() {
	var err error

	// Connect to PostgreSQL
	db, err = sql.Open("pgx", "postgres://user:password@postgres:5432/reelmetrics_db")
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL:", err)
	}

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

// Get all theaters
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

// Get movies & sales for a theater (cached in Redis)
func getMoviesByTheater(w http.ResponseWriter, r *http.Request) {
	theaterID := r.URL.Path[len("/theaters/"):]
	cacheKey := "theater_" + theaterID + "_movies"

	// Check Redis first
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"source": "cache", "data": %s}`, val)))
		return
	}

	// Query PostgreSQL if not cached
	query := `
    SELECT m.id, m.title, SUM(s.tickets_sold * s.ticket_price) AS revenue
    FROM movies m
    JOIN sales s ON m.id = s.movie_id
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
	date := r.URL.Path[len("/top-theater/"):]
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
    SELECT t.name, SUM(s.tickets_sold * s.ticket_price) AS total_revenue
    FROM sales s
    JOIN theaters t ON s.theater_id = t.id
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
	http.HandleFunc("/theaters", getTheaters)
	http.HandleFunc("/theaters/", getMoviesByTheater)
	http.HandleFunc("/top-theater/", getTopTheater)

	log.Println("🚀 Go backend running on port 8002...")
	http.ListenAndServe(":8002", nil)
}
