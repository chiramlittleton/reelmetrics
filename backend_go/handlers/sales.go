package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/db"
	"backend/redis"

	"github.com/gorilla/mux"
)

// GetMoviesByTheater retrieves sales data for a theater from Redis or PostgreSQL
func GetMoviesByTheater(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	theaterID := vars["theater_id"]

	cacheKey := fmt.Sprintf("sales_theater:%s", theaterID)

	// Try to get data from Redis first
	salesList, err := redis.RedisClient.LRange(redis.RedisCtx, cacheKey, 0, -1).Result()
	if err == nil && len(salesList) > 0 {
		log.Printf("✅ Cache hit for %s", cacheKey)
		var salesData []map[string]interface{}

		for _, saleJSON := range salesList {
			var sale map[string]interface{}
			json.Unmarshal([]byte(saleJSON), &sale)
			salesData = append(salesData, sale)
		}

		response, _ := json.Marshal(map[string]interface{}{
			"source": "cache",
			"data":   salesData,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}

	log.Printf("❌ Cache miss for %s, querying PostgreSQL...", cacheKey)

	// If Redis cache is empty, query PostgreSQL
	query := `
		SELECT m.title AS movie_title, s.sale_date, SUM(s.tickets_sold * s.ticket_price) AS revenue
		FROM sales s
		JOIN movies m ON s.movie_id = m.id
		WHERE s.theater_id = $1
		GROUP BY m.title, s.sale_date;
	`

	rows, err := db.DB.Query(query, theaterID)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var title string
		var saleDate string
		var revenue float64
		err := rows.Scan(&title, &saleDate, &revenue)
		if err != nil {
			log.Printf("⚠️ Error scanning row: %v", err)
			continue
		}
		movies = append(movies, map[string]interface{}{
			"title":        title,
			"sale_date":    saleDate,
			"ticket_sales": revenue,
		})
	}

	// Store result in Redis for future queries
	jsonData, _ := json.Marshal(movies)
	redis.RedisClient.RPush(redis.RedisCtx, cacheKey, jsonData)

	// Return the data to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
