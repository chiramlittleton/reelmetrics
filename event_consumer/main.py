import json
import os
import time
import psycopg2
from kafka import KafkaConsumer, KafkaAdminClient
from kafka.admin import NewTopic
from kafka.errors import TopicAlreadyExistsError, NoBrokersAvailable

# Kafka Config
KAFKA_BROKER = os.getenv("KAFKA_BROKER", "kafka:9092")
TOPIC = "ticket_purchases"

# Retry configuration
MAX_RETRIES = 5
RETRY_DELAY = 5  # seconds

def connect_to_postgres():
    """Attempt to establish a connection to PostgreSQL with retries."""
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            conn = psycopg2.connect(
                dbname=os.getenv("POSTGRES_DB", "reelmetrics_db"),
                user=os.getenv("POSTGRES_USER", "user"),
                password=os.getenv("POSTGRES_PASSWORD", "password"),
                host=os.getenv("POSTGRES_HOST", "postgres"),
                port=os.getenv("POSTGRES_PORT", "5432")
            )
            conn.autocommit = True
            print("✅ Connected to PostgreSQL")
            return conn
        except psycopg2.OperationalError as e:
            print(f"⚠️ PostgreSQL connection attempt {attempt} failed: {e}")
            if attempt < MAX_RETRIES:
                time.sleep(RETRY_DELAY)
            else:
                print("❌ Could not connect to PostgreSQL after multiple retries. Exiting.")
                exit(1)

# Initialize DB connection globally
DB_CONN = connect_to_postgres()

def create_kafka_topic():
    """Ensure the Kafka topic exists before consuming messages, with retries."""
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            admin_client = KafkaAdminClient(bootstrap_servers=KAFKA_BROKER)
            topic = NewTopic(name=TOPIC, num_partitions=1, replication_factor=1)

            admin_client.create_topics([topic])
            print(f"✅ Created topic: {TOPIC}")
            return
        except TopicAlreadyExistsError:
            print(f"⚠️ Topic {TOPIC} already exists.")
            return
        except NoBrokersAvailable as e:
            print(f"⚠️ Kafka broker unavailable (attempt {attempt}): {e}")
            if attempt < MAX_RETRIES:
                time.sleep(RETRY_DELAY)
            else:
                print("❌ Could not connect to Kafka after multiple retries. Exiting.")
                exit(1)

def get_movie_id(movie_title):
    """Retrieve the movie_id based on movie_title."""
    with DB_CONN.cursor() as cur:
        cur.execute("SELECT id FROM movies WHERE title = %s;", (movie_title,))
        result = cur.fetchone()
        return result[0] if result else None

def get_theater_id(theater_name):
    """Retrieve the theater_id based on theater_name."""
    with DB_CONN.cursor() as cur:
        cur.execute("SELECT id FROM theaters WHERE name = %s;", (theater_name,))
        result = cur.fetchone()
        return result[0] if result else None

def consume_messages():
    """Consume Kafka messages and write to PostgreSQL with proper lookups."""
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            consumer = KafkaConsumer(
                TOPIC,
                bootstrap_servers=KAFKA_BROKER,
                auto_offset_reset="earliest",
                group_id="sales-consumer-group",
                enable_auto_commit=True,
                value_deserializer=lambda v: json.loads(v.decode("utf-8"))
            )
            print("🚀 Kafka Consumer Started. Listening for messages...")
            break
        except NoBrokersAvailable as e:
            print(f"⚠️ Kafka connection attempt {attempt} failed: {e}")
            if attempt < MAX_RETRIES:
                time.sleep(RETRY_DELAY)
            else:
                print("❌ Could not connect to Kafka after multiple retries. Exiting.")
                exit(1)

    global DB_CONN  # Ensure we can modify the global connection
    for message in consumer:
        event = message.value
        print(f"📥 Received message: {event}")  # Debugging step

        required_keys = {"movie_title", "theater_name", "sale_date", "tickets_sold", "ticket_price"}
        if not required_keys.issubset(event.keys()):
            print(f"⚠️ Skipping message due to missing keys: {event}")
            continue  # Skip this message

        movie_title = event["movie_title"]
        theater_name = event["theater_name"]
        sale_date = event["sale_date"]
        tickets_sold = event["tickets_sold"]
        ticket_price = event["ticket_price"]

        try:
            with DB_CONN.cursor() as cur:
                # Look up movie_id and theater_id
                movie_id = get_movie_id(movie_title)
                theater_id = get_theater_id(theater_name)

                if not movie_id:
                    print(f"⚠️ Movie '{movie_title}' not found in DB. Skipping event.")
                    continue  # Skip if movie is not found
                if not theater_id:
                    print(f"⚠️ Theater '{theater_name}' not found in DB. Skipping event.")
                    continue  # Skip if theater is not found

                # Insert the sale
                cur.execute(
                    """
                    INSERT INTO sales (movie_id, theater_id, sale_date, tickets_sold, ticket_price)
                    VALUES (%s, %s, %s, %s, %s)
                    """,
                    (movie_id, theater_id, sale_date, tickets_sold, ticket_price)
                )
            print(f"✅ Inserted sale: {event}")
        except psycopg2.OperationalError as e:
            print(f"⚠️ Database operation failed: {e}")
            DB_CONN = connect_to_postgres()  # Reconnect and retry

if __name__ == "__main__":
    create_kafka_topic()  # Ensure topic exists
    consume_messages()    # Start consuming messages
