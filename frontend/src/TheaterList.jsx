import React from "react";

function TheaterList({ theaters, fetchSalesForTheater }) {
  return (
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
  );
}

export default TheaterList;
