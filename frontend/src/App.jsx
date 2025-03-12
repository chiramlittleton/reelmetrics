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
  const [selectedDate, setSelectedDate] = useState("");

  useEffect(() => {
    // Clear selections when switching backends
    setSelectedTheater(null);
    setSalesData({});
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

      // ✅ Ensure we correctly extract sales data (supports both Go & Python backends)
      const salesArray = Array.isArray(response.data?.data) ? response.data.data : response.data;

      const salesByDate = {};
      salesArray.forEach((sale) => {
        const saleDate = sale.sale_date.split("T")[0]; // Normalize date format
        if (!salesByDate[saleDate]) {
          salesByDate[saleDate] = [];
        }
        salesByDate[saleDate].push({
          title: sale.movie_title, // ✅ Ensure correct field name
          tickets_sold: sale.tickets_sold, // ✅ Store ticket count
          ticket_price: sale.ticket_price, // ✅ Store ticket price
          revenue: sale.tickets_sold * sale.ticket_price, // ✅ Compute revenue
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

      if (response.data?.theater && response.data?.revenue !== undefined) {
        return response.data; // ✅ Correctly return structured response
      }

      return { message: "No sales data available" }; // ✅ Handle missing data case
    } catch (error) {
      console.error("Error fetching top theater", error);
      return { message: "Error fetching data" };
    }
  };

  return (
    <div>
      <h1>ReelMetrics</h1>

      <BackendSelector backend={backend} setBackend={setBackend} />
      <TheaterList theaters={theaters} fetchSalesForTheater={fetchSalesForTheater} />
      <SalesData selectedTheater={selectedTheater} salesData={salesData} />
      <TopTheater fetchTopTheater={fetchTopTheater} selectedDate={selectedDate} />
    </div>
  );
}

export default App;
