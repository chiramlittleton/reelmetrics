FROM python:3.10

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application files
COPY . .

# Set environment variables
ENV POSTGRES_HOST="postgres"
ENV REDIS_HOST="redis"
ENV KAFKA_BROKER="kafka:9092"

# Run the event consumer
CMD ["python", "main.py"]
