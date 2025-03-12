from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import psycopg2
import redis
import json
import os
from datetime import date
from decimal import Decimal

app = FastAPI()

# ✅ Add CORS Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Allow all origins (or restrict to ["http://localhost:3000"])
    allow_credentials=True,
    allow_methods=["*"],  # Allow all HTTP methods (GET, POST, etc.)
    allow_headers=["*"],  # Allow all headers
)

# PostgreSQL Connection
DB_CONN = psycopg2.connect(
    dbname=os.getenv("POSTGRES_DB", "reelmetrics_db"),
    user=os.getenv("POSTGRES_USER", "user"),
    password=os.getenv("POSTGRES_PASSWORD", "password"),
    host=os.getenv("POSTGRES_HOST", "postgres"),
    port=os.getenv("POSTGRES_PORT", "5432")
)
DB_CONN.autocommit = True

# Redis connection
redis_client = redis.Redis(host=os.getenv("REDIS_HOST", "redis"), port=6379, decode_responses=True)

@app.get("/theaters")
def get_theaters():
    """Fetch all theaters from PostgreSQL"""
    with DB_CONN.cursor() as cur:
        cur.execute("SELECT id, name FROM theaters;")
        theaters = [{"id": row[0], "name": row[1]} for row in cur.fetchall()]
    return theaters

@app.get("/theaters/{theater_id}/movies")
def get_movies_by_theater(theater_id: int):
    """Fetch movies & sales for a theater, using Redis cache if available"""
    cache_key = f"theater:{theater_id}:movies"
    cached_data = redis_client.get(cache_key)

    if cached_data:
        return json.loads(cached_data)  # Return cached data

    with DB_CONN.cursor() as cur:
        cur.execute(
            """
            SELECT m.id, m.title, COALESCE(SUM(s.tickets_sold * s.ticket_price), 0) AS revenue
            FROM movies m
            LEFT JOIN sales s ON m.id = s.movie_id
            WHERE m.theater_id = %s
            GROUP BY m.id, m.title;
            """,
            (theater_id,)
        )
        movies = [{"id": row[0], "title": row[1], "ticket_sales": float(row[2])} for row in cur.fetchall()]

    redis_client.setex(cache_key, 60, json.dumps(movies))  # Cache for 60 seconds
    return movies

@app.get("/top-theater/{sale_date}")
def get_top_theater(sale_date: date):
    """Fetch the top theater by revenue for a given date"""
    cache_key = f"top_theater:{sale_date}"
    cached_data = redis_client.get(cache_key)

    if cached_data:
        return json.loads(cached_data)  # Return cached data

    with DB_CONN.cursor() as cur:
        cur.execute(
            """
            SELECT t.name, COALESCE(SUM(s.tickets_sold * s.ticket_price), 0) AS total_revenue
            FROM sales s
            JOIN movies m ON s.movie_id = m.id
            JOIN theaters t ON m.theater_id = t.id
            WHERE s.sale_date = %s
            GROUP BY t.name
            ORDER BY total_revenue DESC
            LIMIT 1;
            """,
            (sale_date,)
        )
        result = cur.fetchone()
        if not result:
            raise HTTPException(status_code=404, detail="No sales data for this date")

        top_theater = {"theater": result[0], "revenue": float(result[1])}
        redis_client.setex(cache_key, 300, json.dumps(top_theater))  # Cache for 5 minutes

    return top_theater
