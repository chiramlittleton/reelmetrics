from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import psycopg2
import redis
import json
import os
from decimal import Decimal
from typing import List, Dict

app = FastAPI()

# Add CORS Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# PostgreSQL Connection
DATABASE_URL = os.getenv("DATABASE_URL", "postgres://user:password@postgres:5432/reelmetrics_db")
DB_CONN = psycopg2.connect(DATABASE_URL)

# Redis Connection
redis_client = redis.Redis(host=os.getenv("REDIS_HOST", "redis"), port=6379, decode_responses=True)

def decimal_to_float(obj):
    """Convert Decimal objects to float for JSON serialization."""
    if isinstance(obj, Decimal):
        return float(obj)
    raise TypeError(f"Type {type(obj)} not serializable")

@app.get("/theaters", response_model=List[Dict])
def get_theaters():
    """Fetch all theaters from PostgreSQL"""
    with DB_CONN.cursor() as cur:
        cur.execute("SELECT id, name FROM theaters;")
        theaters = [{"id": row[0], "name": row[1]} for row in cur.fetchall()]
    return theaters

@app.get("/theaters/{theater_id}/movies")
def get_movies_by_theater(theater_id: int):
    """Fetch movies & sales for a theater, using Redis cache if available"""

    redis_key = f"sales_theater:{theater_id}"
    
    # Check Redis cache first
    sales_list = redis_client.lrange(redis_key, 0, -1)
    if sales_list:
        print(f"✅ Cache hit for {redis_key}")
        return [json.loads(sale) for sale in sales_list]

    print(f"❌ Cache miss for {redis_key}, querying PostgreSQL...")

    # Fetch from PostgreSQL if not cached
    with DB_CONN.cursor() as cur:
        cur.execute(
            """
            SELECT m.id, m.title, s.sale_date, SUM(s.tickets_sold * s.ticket_price) AS revenue
            FROM movies m
            JOIN sales s ON m.id = s.movie_id
            WHERE s.theater_id = %s
            GROUP BY m.id, m.title, s.sale_date;
            """,
            (theater_id,)
        )
        movies = [
            {"id": row[0], "title": row[1], "sale_date": str(row[2]), "ticket_sales": float(row[3])}
            for row in cur.fetchall()
        ]

    # Store result in Redis as a list
    for movie in movies:
        redis_client.rpush(redis_key, json.dumps(movie, default=decimal_to_float))

    return movies

@app.get("/top-theater/{sale_date}")
def get_top_theater(sale_date: str):
    """Fetch the top theater by sales on a given date, using Redis"""

    redis_key = f"sales_date:{sale_date}"
    sales_list = redis_client.lrange(redis_key, 0, -1)

    if not sales_list:
        print(f"❌ No sales data in Redis for {sale_date}")
        return {"message": "No sales data available"}

    print(f"✅ Found sales data in Redis for {sale_date}")

    # Compute the top theater by revenue
    theater_revenue = {}

    for sale_json in sales_list:
        sale = json.loads(sale_json)
        theater_id = sale["theater_id"]
        revenue = sale["tickets_sold"] * sale["ticket_price"]

        theater_revenue[theater_id] = theater_revenue.get(theater_id, 0) + revenue

    # Find the highest revenue theater
    top_theater_id = max(theater_revenue, key=theater_revenue.get, default=None)

    if not top_theater_id:
        return {"message": "No sales data available"}

    # Retrieve theater name (check Redis first, then PostgreSQL)
    redis_theater_key = f"theater:{top_theater_id}"
    theater_name = redis_client.get(redis_theater_key)

    if not theater_name:
        with DB_CONN.cursor() as cur:
            cur.execute("SELECT name FROM theaters WHERE id = %s;", (top_theater_id,))
            result = cur.fetchone()
            if result:
                theater_name = result[0]
                redis_client.set(redis_theater_key, theater_name)  # Cache it

    return {"theater": theater_name, "revenue": theater_revenue[top_theater_id]}
