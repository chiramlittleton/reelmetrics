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
    with DB_CONN.cursor() as cur:
        cur.execute(
            """
            SELECT m.id, m.title, s.sale_date, SUM(s.tickets_sold * s.ticket_price) AS revenue
            FROM movies m
            JOIN sales s ON m.id = s.movie_id
            WHERE m.theater_id = %s
            GROUP BY m.id, m.title, s.sale_date;
            """,
            (theater_id,)
        )
        movies = [{"id": row[0], "title": row[1], "sale_date": str(row[2]), "ticket_sales": float(row[3])} for row in cur.fetchall()]
    return movies