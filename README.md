# ReelMetrics

**ReelMetrics** is a real-time ticket sales tracking system designed for movie theaters. It captures and processes ticket purchases as they happen, ensuring up-to-date sales data and analytics. The system is built with a **modular, microservices-based architecture**, using **Kafka** for seamless event communication, **PostgreSQL** for reliable data storage, and **Redis** for quick lookups. With both **Python (FastAPI)** and **Go** backends, ReelMetrics showcases a flexible, multi-language approach while delivering fast and accurate insights into theater performance and movie sales.


## Architecture Overview
ReelMetrics follows a microservices-based, event-driven design, where Kafka streams ticket sales to a consumer that updates PostgreSQL and Redis. The Go and Python (FastAPI) backends query both databases, exposing APIs for theater and sales data. A React frontend dynamically switches between backends and fetches analytics.

By separating event ingestion, storage, caching, and APIs, the system ensures efficient data flow, low-latency queries, and scalability.

```mermaid
graph TD
    A("Event Generator") -- Ticket Sales --> B("Kafka")
    B -- Events --> C("Event Consumer")
    C -- Updates --> D["Postgres"] & E["Redis"]
    D -- Queries --- F("Backend (Go)")
    E -- Queries --- F
    %% F -- Init --> D
    D -- Queries --- G("Backend (Python)")
    E -- Queries --- G
    F -- REST --- H("Frontend (React)")
    G -- REST --- H("Frontend (React)")


    D@{ shape: cyl}
    E@{ shape: stored-data}
    style B fill:#FF6D00
    style E color:#FFFFFF, fill:#AA00FF, stroke:#AA00FF

```

## 🛠️ Tech Stack

  - **Backend (Python - FastAPI)**: Exposes APIs to fetch theaters & movie sales.
  - **Backend (Go)**: Alternative implementation for fetching the same data.
  - **PostgreSQL**: Stores theaters, movies, and sales data.
  - **Redis**: Caches frequently accessed data (e.g., top theaters, movie sales).
  - **Kafka**: Streams real-time ticket purchase events.
  - **Frontend (React)**: Displays theaters, movies, and sales statistics.

## **🧩 Project Components**

### **Event Generator**
- Generates **random ticket sales events**.
- Publishes events to **Kafka**.

### **Backends (Python & Go)**
- **Fetch theaters & movie sales from PostgreSQL**.
- **Cache results in Redis** to optimize performance.

### **Frontend (React)**
- Lets users select between **Python & Go backends**.
- Displays **real-time revenue statistics**.

## **📦 Setup & Installation**

### **1️⃣ Clone the Repository**
```bash
git clone https://github.com/yourusername/reelmetrics.git
cd reelmetrics
```

### **2️⃣ Start the Services with Docker**
```bash
docker-compose up -d
```

### **3️⃣ Verify Running Services**
```bash
docker ps
```

### **4️⃣ Test API Endpoints**

#### **Get Theaters (Python Backend)**
```bash
curl -X GET http://localhost:8001/theaters
```

#### **Get Movies & Sales for a Theater (`id=1`)**
```bash
curl -X GET http://localhost:8001/theaters/1/movies
```

#### **Get Top Theater by Revenue (Go Backend)**
```bash
curl -X GET http://localhost:8002/top-theater/2024-05-10
```

## **🔧 Development**
### **Run Python Backend Locally**
```bash
cd backend_python
uvicorn main:app --host 0.0.0.0 --port 8001 --reload
```

### **Run Go Backend Locally**
```bash
cd backend_go
go run main.go
```

### **Run Event Generator**
```bash
cd utils
python event_simulator.py
```

## **📌 Next Steps**
- Add **user authentication** to restrict access.
- Implement **real-time WebSockets** for live sales updates.
- Deploy on **AWS using Kubernetes**.

## **📝 License**
This project is licensed under the MIT License.


### Screenshots

Below are some screenshots of the application in action.

#### 🎭 Theaters List
![Theaters List](./screenshots/theaters_list.png)

#### 🏆 Top Theater by Revenue
![Top Theater](./screenshots/top_theater.png)


