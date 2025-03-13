package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"backend/db"
	"backend/redis"

	"github.com/gorilla/mux"
)

// GetTopTheater fetches the highest revenue theater for a given date
func GetTopTheater(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]

	redisKey := "sales_date:" + date

	// Fetch sales from Redis using LRANGE
	salesList, err := redis.RedisClient.LRange(redis.RedisCtx, redisKey, 0, -1).Result()
	if err != nil || len(salesList) == 0 {
		log.Printf("❌ No sales data found in Redis for date %s", date)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "No sales data available"}`))
		return
	}

	log.Printf("✅ Found sales data in Redis for %s", date)

	// Process sales data to calculate top revenue
	theaterRevenue := make(map[int]float64) // Store revenue by theater_id

	for _, saleJSON := range salesList {
		var sale map[string]interface{}
		json.Unmarshal([]byte(saleJSON), &sale)

		theaterID, err := strconv.Atoi(fmt.Sprintf("%v", sale["theater_id"]))
		if err != nil {
			log.Println("⚠️ Invalid theater_id in sale data")
			continue
		}

		ticketSales, err := strconv.ParseFloat(fmt.Sprintf("%v", sale["tickets_sold"]), 64)
		if err != nil {
			log.Println("⚠️ Invalid tickets_sold value")
			continue
		}

		ticketPrice, err := strconv.ParseFloat(fmt.Sprintf("%v", sale["ticket_price"]), 64)
		if err != nil {
			log.Println("⚠️ Invalid ticket_price value")
			continue
		}

		revenue := ticketSales * ticketPrice
		theaterRevenue[theaterID] += revenue
	}

	// Find the top theater by revenue
	var topTheaterID int
	var maxRevenue float64

	for theaterID, revenue := range theaterRevenue {
		if revenue > maxRevenue {
			maxRevenue = revenue
			topTheaterID = theaterID
		}
	}

	// If no valid sales exist, return a "no data" message
	if topTheaterID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "No sales data available"}`))
		return
	}

	// Fetch the theater name (either from Redis or PostgreSQL)
	theaterName, err := getTheaterName(topTheaterID)
	if err != nil {
		log.Printf("❌ Failed to fetch theater name for ID %d", topTheaterID)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Error retrieving theater name"}`))
		return
	}

	// Return the top theater data
	result := map[string]interface{}{
		"theater": theaterName,
		"revenue": maxRevenue,
	}
	jsonData, _ := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// getTheaterName retrieves the theater name given its ID
func getTheaterName(theaterID int) (string, error) {
	// First, check if the theater name is stored in Redis
	redisKey := fmt.Sprintf("theater:%d", theaterID)
	theaterName, err := redis.RedisClient.Get(redis.RedisCtx, redisKey).Result()
	if err == nil {
		return theaterName, nil // Return from Redis if available
	}

	// If not found in Redis, query PostgreSQL
	query := "SELECT name FROM theaters WHERE id = $1;"
	row := db.DB.QueryRow(query, theaterID)
	err = row.Scan(&theaterName)
	if err != nil {
		return "", err
	}

	// Store the result in Redis for future lookups
	redis.RedisClient.Set(redis.RedisCtx, redisKey, theaterName, 0)
	return theaterName, nil
}
