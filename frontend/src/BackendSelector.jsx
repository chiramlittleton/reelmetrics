import React from "react";

const BACKENDS = {
  python: "http://localhost:8001",
  go: "http://localhost:8002",
};

function BackendSelector({ backend, setBackend }) {
  return (
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
  );
}

export default BackendSelector;
