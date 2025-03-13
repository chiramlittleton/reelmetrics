# ReelMetrics

**ReelMetrics** is a real-time ticket sales tracking system designed for movie theaters. It captures and processes ticket purchases as they happen, ensuring up-to-date sales data and analytics. The system is built with a **modular, microservices-based architecture**, using **Kafka** for seamless event communication, **PostgreSQL** for reliable data storage, and **Redis** for quick lookups. With both **Python (FastAPI)** and **Go** backends, ReelMetrics showcases a flexible, multi-language approach while delivering fast and accurate insights into theater performance and movie sales.


## Architecture Overview
ReelMetrics follows a microservices-based, event-driven design, where Kafka streams ticket sales to a consumer that updates PostgreSQL and Redis. The Go and Python (FastAPI) backends query both databases, exposing APIs for theater and sales data. A React frontend dynamically switches between backends and fetches analytics.

By separating event ingestion, storage, caching, and APIs, the system ensures efficient data flow, low-latency queries, and scalability.

```mermaid
flowchart TD
%% Nodes
    A("fab:fa-youtube Starter Guide")
    B("fab:fa-youtube Make Flowchart")
    n1@{ icon: "fa:gem", pos: "b", h: 24}
    C("fa:fa-book-open Learn More")
    D{"Use the editor"}
    n2(Many shapes)@{ shape: delay}
    E(fa:fa-shapes Visual Editor)
    F("fa:fa-chevron-up Add node in toolbar")
    G("fa:fa-comment-dots AI chat")
    H("fa:fa-arrow-left Open AI in side menu")
    I("fa:fa-code Text")
    J(fa:fa-arrow-left Type Mermaid syntax)

%% Edge connections between nodes
    A --> B --> C --> n1 & D & n2
    D -- Build and Design --> E --> F
    D -- Use AI --> G --> H
    D -- Mermaid js --> I --> J

%% Individual node styling. Try the visual editor toolbar for easier styling!
    style E color:#FFFFFF, fill:#AA00FF, stroke:#AA00FF
    style G color:#FFFFFF, stroke:#00C853, fill:#00C853
    style I color:#FFFFFF, stroke:#2962FF, fill:#2962FF

%% You can add notes with two "%" signs in a row!
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


