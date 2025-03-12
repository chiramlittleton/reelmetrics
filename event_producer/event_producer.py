import json
import random
import time
from kafka import KafkaProducer

KAFKA_BROKER = "127.0.0.1:9093"
TOPIC = "ticket_purchases"

THEATERS = [
    "AMC Century City",
    "Regal LA Live",
    "Cinemark Playa Vista"
]

MOVIES = [
    "Jurassic Park",
    "Titanic",
    "The Matrix",
    "Pulp Fiction",
    "The Lion King",
    "Forrest Gump"
]

DATES = ["2025-03-11", "2025-03-12", "2025-03-13"]

TICKET_PRICE_RANGE = (5.00, 15.00)

producer = KafkaProducer(
    bootstrap_servers=KAFKA_BROKER,
    value_serializer=lambda v: json.dumps(v).encode("utf-8")
)

def generate_ticket_purchase():
    """Generate a random ticket purchase event."""
    theater_name = random.choice(THEATERS)  # Renamed to match database schema
    movie_title = random.choice(MOVIES)  # Renamed to match database schema
    sale_date = random.choice(DATES)
    ticket_price = round(random.uniform(*TICKET_PRICE_RANGE), 2)  # Random ticket price between $5 - $15
    tickets_sold = random.randint(1, 5)  # Random number of tickets per transaction (1-5)

    return {
        "theater_name": theater_name,  # Match database column name
        "movie_title": movie_title,  # Match database column name
        "sale_date": sale_date,
        "ticket_price": ticket_price,
        "tickets_sold": tickets_sold
    }

if __name__ == "__main__":
    print("Starting ticket purchase simulation...")
    try:
        while True:
            ticket_event = generate_ticket_purchase()
            producer.send(TOPIC, ticket_event)
            print(f"Sent event: {ticket_event}")
            time.sleep(2)  # Generate a new ticket purchase every 2 seconds
    except KeyboardInterrupt:
        print("Simulation stopped.")
