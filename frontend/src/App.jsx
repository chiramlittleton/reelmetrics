import React, { useState, useEffect } from "react";
import axios from "axios";

const BACKENDS = {
  python: "http://localhost:8001",
  go: "http://localhost:8002"
};

function App() {
  const [backend, setBackend] = useState("python");
  const [theaters, setTheaters] = useState([]);
  const [selectedTheater, setSelectedTheater] = useState(null);
  const [salesData, setSalesData] = useState({});
  const [topTheater, setTopTheater] = useState(null);

  useEffect(() => {
    fetchTheaters();
  }, [backend]);

  const fetchTheaters = async () => {
    try {
      const response = await axios.get(`${BACKENDS[backend]}/theaters`);
      setTheaters(response.data);
    } catch (error) {
      console.error("Error fetching theaters", error);
    }
  };

  const fetchSalesForTheater = async (theaterId) => {
    try {
      setSelectedTheater(theaterId);
      const response = await axios.get(`${BACKENDS[backend]}/theaters/${theaterId}/movies`);
      
      const salesByDate = {};
      response.data.forEach((sale) => {
        if (!salesByDate[sale.sale_date]) {
          salesByDate[sale.sale_date] = [];
        }
        salesByDate[sale.sale_date].push({
          title: sale.title,
          ticket_sales: sale.ticket_sales
        });
      });

      setSalesData(salesByDate);
    } catch (error) {
      console.error("Error fetching sales", error);
    }
  };

  const fetchTopTheater = async (saleDate) => {
    try {
      const response = await axios.get(`${BACKENDS[backend]}/top-theater/${saleDate}`);
      setTopTheater(response.data);
    } catch (error) {
      console.error("Error fetching top theater", error);
    }
  };

  return (
    <div>
      <h1>ReelMetrics</h1>

      <div>
        <h2>Select Backend:</h2>
        {Object.keys(BACKENDS).map((key) => (
          <label key={key}>
            <input
              type="radio"
              value={key}
              checked={backend === key}
              onChange={() => setBackend(key)}
            />
            {key.toUpperCase()}
          </label>
        ))}
      </div>

      <div>
        <h2>Select a Theater:</h2>
        {theaters.length > 0 ? (
          <ul>
            {theaters.map((theater) => (
              <li key={theater.id} onClick={() => fetchSalesForTheater(theater.id)}>
                {theater.name}
              </li>
            ))}
          </ul>
        ) : (
          <p>No theaters available.</p>
        )}
      </div>

      {selectedTheater && (
        <div>
          <h2>Sales Data</h2>
          {Object.keys(salesData).length > 0 ? (
            Object.entries(salesData).map(([date, movies]) => (
              <div key={date}>
                <h3>{date}</h3>
                <ul>
                  {movies.map((movie, index) => (
                    <li key={index}>
                      {movie.title} - ${movie.ticket_sales.toFixed(2)}
                    </li>
                  ))}
                </ul>
              </div>
            ))
          ) : (
            <p>No sales data found for this theater.</p>
          )}
        </div>
      )}

      <div>
        <h2>Top Theater by Sales</h2>
        <input
          type="date"
          onChange={(e) => fetchTopTheater(e.target.value)}
        />
        {topTheater ? (
          <p>
            <strong>{topTheater.theater}</strong> - ${topTheater.revenue.toFixed(2)}
          </p>
        ) : (
          <p>Select a date to see the top theater.</p>
        )}
      </div>
    </div>
  );
}

export default App;
