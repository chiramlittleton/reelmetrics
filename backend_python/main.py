from fastapi import FastAPI, HTTPException
import psycopg2
import redis
import json
import os

app = FastAPI()

# PostgreSQL connection
DB_CONN = psycopg2.connect(
    dbname=os.getenv("POSTGRES_DB", "reelmetrics"),
    user=os.getenv("POSTGRES_USER", "postgres"),
    password=os.getenv("POSTGRES_PASSWORD", "password"),
    host=os.getenv("POSTGRES_HOST", "postgres"),
    port=os.getenv("POSTGRES_PORT", "5432")
)

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

    # Fetch from PostgreSQL if not cached
    with DB_CONN.cursor() as cur:
        cur.execute(
            """
            SELECT m.id, m.title, SUM(s.tickets_sold * s.ticket_price) AS revenue
            FROM movies m
            JOIN sales s ON m.id = s.movie_id
            WHERE m.theater_id = %s
            GROUP BY m.id, m.title;
            """,
            (theater_id,)
        )
        movies = [{"id": row[0], "title": row[1], "ticket_sales": row[2]} for row in cur.fetchall()]

    # Cache the result in Redis
    redis_client.setex(cache_key, 60, json.dumps(movies))  # Cache for 60 seconds

    return movies
