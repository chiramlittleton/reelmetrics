import React, { useState, useEffect } from "react";
import axios from "axios";
import BackendSelector from "./BackendSelector";
import TheaterList from "./TheaterList";
import SalesData from "./SalesData";
import TopTheater from "./TopTheater";

const BACKENDS = {
  python: "http://localhost:8001",
  go: "http://localhost:8002",
};

function App() {
  const [backend, setBackend] = useState("python");
  const [theaters, setTheaters] = useState([]);
  const [selectedTheater, setSelectedTheater] = useState(null);
  const [salesData, setSalesData] = useState({});
  const [topTheater, setTopTheater] = useState(null);
  const [selectedDate, setSelectedDate] = useState("");

  useEffect(() => {
    // Clear selections when switching backends
    setSelectedTheater(null);
    setSalesData({});
    setTopTheater(null);
    setSelectedDate(""); // Clear date selection
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

      const salesArray = response.data?.data || response.data || [];

      const salesByDate = {};
      salesArray.forEach((sale) => {
        if (!salesByDate[sale.sale_date]) {
          salesByDate[sale.sale_date] = [];
        }
        salesByDate[sale.sale_date].push({
          title: sale.title,
          ticket_sales: sale.ticket_sales,
        });
      });

      setSalesData(salesByDate);
    } catch (error) {
      console.error("Error fetching sales", error);
    }
  };

  const fetchTopTheater = async (saleDate) => {
    try {
      setSelectedDate(saleDate); // Update selected date
      const response = await axios.get(`${BACKENDS[backend]}/top-theater/${saleDate}`);
      setTopTheater(response.data);
    } catch (error) {
      console.error("Error fetching top theater", error);
    }
  };

  return (
    <div>
      <h1>ReelMetrics</h1>

      <BackendSelector backend={backend} setBackend={setBackend} />
      <TheaterList theaters={theaters} fetchSalesForTheater={fetchSalesForTheater} />
      <SalesData selectedTheater={selectedTheater} salesData={salesData} />
      <TopTheater 
        fetchTopTheater={fetchTopTheater} 
        topTheater={topTheater} 
        selectedDate={selectedDate} 
      />
    </div>
  );
}

export default App;
