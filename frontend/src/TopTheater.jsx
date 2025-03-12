import React, { useState, useEffect } from "react";

function TopTheater({ fetchTopTheater, selectedDate }) {
  const [topTheater, setTopTheater] = useState(null);
  const [noDataMessage, setNoDataMessage] = useState(""); // Store "no data" message

  useEffect(() => {
    const getTopTheater = async () => {
      if (!selectedDate) {
        setTopTheater(null);
        setNoDataMessage("");
        return;
      }
      const data = await fetchTopTheater(selectedDate);

      if (data?.message) {
        setNoDataMessage(data.message);
        setTopTheater(null);
      } else {
        setNoDataMessage("");
        setTopTheater(data);
      }
    };

    getTopTheater();
  }, [selectedDate, fetchTopTheater]);

  return (
    <div>
      <h2>Top Theater by Sales</h2>
      <input
        type="date"
        value={selectedDate}
        onChange={(e) => fetchTopTheater(e.target.value)}
      />
      {noDataMessage ? (
        <p>{noDataMessage}</p>
      ) : topTheater ? (
        <p>
          <strong>{topTheater.theater}</strong> - ${topTheater.revenue.toFixed(2)}
        </p>
      ) : (
        <p>Select a date to see the top theater.</p>
      )}
    </div>
  );
}

export default TopTheater;
