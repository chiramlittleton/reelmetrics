package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backend/db"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	log.Println("✅ Connected to Redis")
}

// ✅ Initializes Redis with existing sales data
func InitializeRedisCache() {
	log.Println("🔄 Initializing Redis with existing sales data...")

	query := `
		SELECT s.id, s.theater_id, m.title AS movie_title, s.sale_date, s.tickets_sold, s.ticket_price
		FROM sales s
		JOIN movies m ON s.movie_id = m.id;
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Fatalf("❌ Failed to fetch sales data from PostgreSQL: %v", err)
	}
	defer rows.Close()

	// Process each sale record
	for rows.Next() {
		var saleID int
		var theaterID int
		var movieTitle string
		var rawSaleDate time.Time
		var ticketsSold int
		var ticketPrice float64

		err := rows.Scan(&saleID, &theaterID, &movieTitle, &rawSaleDate, &ticketsSold, &ticketPrice)
		if err != nil {
			log.Printf("⚠️ Skipping record due to scan error: %v", err)
			continue
		}

		saleDate := rawSaleDate.Format("2006-01-02")

		// Create sale data object
		saleData := map[string]interface{}{
			"sale_id":      saleID,
			"theater_id":   theaterID,
			"movie_title":  movieTitle,
			"sale_date":    saleDate,
			"tickets_sold": ticketsSold,
			"ticket_price": ticketPrice,
		}

		saleJSON, _ := json.Marshal(saleData)

		// ✅ Store individual sale (if missing)
		redisKeySale := fmt.Sprintf("sale:%d", saleID)
		exists, _ := RedisClient.Exists(RedisCtx, redisKeySale).Result()
		if exists == 0 {
			RedisClient.Set(RedisCtx, redisKeySale, saleJSON, 0)
		}

		// ✅ Append to sales by theater (ensure no duplicates)
		redisKeyTheater := fmt.Sprintf("sales_theater:%d", theaterID)
		RedisClient.LRem(RedisCtx, redisKeyTheater, 0, saleJSON) // Remove if exists
		RedisClient.RPush(RedisCtx, redisKeyTheater, saleJSON)

		// ✅ Append to sales by date (ensure no duplicates)
		redisKeyDate := fmt.Sprintf("sales_date:%s", saleDate)
		RedisClient.LRem(RedisCtx, redisKeyDate, 0, saleJSON) // Remove if exists
		RedisClient.RPush(RedisCtx, redisKeyDate, saleJSON)
	}

	log.Println("✅ Redis sales cache initialized.")
}
