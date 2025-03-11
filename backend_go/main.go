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
	// Connect to PostgreSQL
	var err error
	db, err = sql.Open("pgx", "postgres://user:password@postgres:5432/reelmetrics_db")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

func getTopTheater(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Path[len("/top-theater/"):]
	cacheKey := "top_theater_" + date

	// Check Redis
	val, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"source": "cache", "data": %s}`, val)))
		return
	}

	// Query PostgreSQL
	query := `
    SELECT t.name, SUM(s.revenue) AS total_revenue
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

	// Store in Redis
	result := map[string]interface{}{"theater": theater, "revenue": revenue}
	jsonData, _ := json.Marshal(result)
	redisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"source": "db", "data": %s}`, jsonData)))
}

func main() {
	http.HandleFunc("/top-theater/", getTopTheater)
	log.Println("Go backend running on port 8002...")
	http.ListenAndServe(":8002", nil)
}
