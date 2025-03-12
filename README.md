# ReelMetrics

ReelMetrics is a distributed system that processes theater ticket sales. It leverages **Kafka** for event streaming, **PostgreSQL** for persistent storage, and **Redis** for caching frequently accessed data. The system includes both **Python (FastAPI)** and **Go** backends to demonstrate multi-language interoperability.

## **🚀 Architecture Overview**
```mermaid
graph TD;
    A[Frontend (React)] -->|API Requests| B[Backend (FastAPI & Go)];
    B -->|Stores Data| C[PostgreSQL];
    B -->|Caches Data| D[Redis];
    E[Event Generator] -->|Publishes Events| F[Kafka];
    F -->|Consumes Events| G[Event Consumer];
    G -->|Writes Data| C;
    G -->|Updates Cache| D;
```

> **Note:** The intended design included an event generator and event consumer to simulate Kafka messages, but this was not fully implemented. The Python backend does not currently use Redis caching.

## **🛠️ Tech Stack**
- **Backend (Python - FastAPI)**: Exposes APIs to fetch theaters & movie sales.
- **Backend (Go)**: Alternative implementation for fetching the same data.
- **PostgreSQL**: Stores theaters, movies, and sales data.
- **Redis**: Caches frequently accessed data (Go backend only).
- **Kafka**: Intended for event streaming but not currently in use.
- **Frontend (React)**: Displays theaters, movies, and sales statistics.

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

#### **Get Sales Data Per Day for a Theater (`id=1`)**
```bash
curl -X GET http://localhost:8001/theaters/1/movies
```

#### **Get Top Theater by Revenue Per Day (Go Backend)**
```bash
curl -X GET http://localhost:8002/top-theater/2025-03-11
```

## **🎬 Project Components**

### **1️⃣ Event Generator (Not Implemented)**
- Originally planned to generate **random ticket sales events**.
- Would publish events to **Kafka**.

### **2️⃣ Backends (Python & Go)**
- **Fetch theaters & sales data from PostgreSQL**.
- **Cache results in Redis (Go only)** to optimize performance.

### **3️⃣ Frontend (React)**
- Lets users select between **Python & Go backends**.
- Displays **daily revenue statistics**.
- Shows **overall top theater per day**.

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

### **Run Event Generator (If Implemented in the Future)**
```bash
cd utils
python event_simulator.py
```

## **🚀 Next Steps**
- Fix **Redis caching for the Python backend**.
- Complete **Kafka event streaming integration**.
- Implement **real-time WebSockets** for live updates.
- Deploy on **AWS using Kubernetes**.

## **📝 License**
This project is licensed under the MIT License.

