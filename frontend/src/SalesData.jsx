import React from "react";

function SalesData({ selectedTheater, salesData }) {
  return (
    <>
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
                      {movie.title} - {movie.tickets_sold} tickets @ ${movie.ticket_price.toFixed(2)} 
                      {" "} = <strong>${movie.revenue.toFixed(2)}</strong>
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
    </>
  );
}

export default SalesData;
