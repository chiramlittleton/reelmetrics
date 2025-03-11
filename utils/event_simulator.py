import json
import random
import time
from kafka import KafkaProducer

KAFKA_BROKER = "localhost:9092"
TOPIC = "ticket_purchases"

THEATERS = [
    "Regal Cinemas",
    "AMC Theaters",
    "Cinemark",
    "Edwards Theaters",
    "Loews Cineplex"
]

MOVIES = [
    "Jurassic Park",
    "Titanic",
    "The Matrix",
    "Pulp Fiction",
    "The Lion King",
    "Forrest Gump",
    "Home Alone",
    "Independence Day",
    "Terminator 2: Judgment Day",
    "Toy Story"
]

DATES = ["1995-07-14", "1996-11-22", "1997-12-19", "1998-07-24", "1999-03-31"]

TICKET_PRICE_RANGE = (5.00, 15.00)

producer = KafkaProducer(
    bootstrap_servers=KAFKA_BROKER,
    value_serializer=lambda v: json.dumps(v).encode("utf-8")
)

def generate_ticket_purchase():
    """Generate a random ticket purchase event."""
    theater = random.choice(THEATERS)
    movie = random.choice(MOVIES)
    sale_date = random.choice(DATES)
    ticket_price = round(random.uniform(*TICKET_PRICE_RANGE), 2)  # Random ticket price between $5 - $15
    tickets_sold = random.randint(1, 5)  # Random number of tickets per transaction (1-5)

    return {
        "theater": theater,
        "movie": movie,
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
