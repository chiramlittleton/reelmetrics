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
  const [movies, setMovies] = useState([]);

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

  const fetchMoviesForTheater = async (theaterId) => {
    try {
      setSelectedTheater(theaterId);
      const response = await axios.get(`${BACKENDS[backend]}/theaters/${theaterId}/movies`);
      setMovies(response.data);
    } catch (error) {
      console.error("Error fetching movies", error);
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
              <li key={theater.id} onClick={() => fetchMoviesForTheater(theater.id)}>
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
          <h2>Movies & Ticket Sales</h2>
          {movies.length > 0 ? (
            <ul>
              {movies.map((movie) => (
                <li key={movie.id}>
                  {movie.title} - ${movie.ticket_sales.toFixed(2)}
                </li>
              ))}
              <li>
                <strong>
                  Total Sales: ${movies.reduce((sum, m) => sum + m.ticket_sales, 0).toFixed(2)}
                </strong>
              </li>
            </ul>
          ) : (
            <p>No movies found for this theater.</p>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
