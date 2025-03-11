# Use official Python image
FROM python:3.10

# Set working directory
WORKDIR /app

# Copy files
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

# Expose FastAPI port
EXPOSE 8001

# Run FastAPI app
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8001"]
