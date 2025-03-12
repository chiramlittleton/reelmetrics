import React from "react";

function TopTheater({ fetchTopTheater, topTheater }) {
  return (
    <div>
      <h2>Top Theater by Sales</h2>
      <input type="date" onChange={(e) => fetchTopTheater(e.target.value)} />
      {topTheater ? (
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
