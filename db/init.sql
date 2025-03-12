-- Drop tables if they exist (optional cleanup for fresh starts)
DROP TABLE IF EXISTS sales, movies, theaters CASCADE;

-- Create Theaters Table
CREATE TABLE theaters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Create Movies Table
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE
);

-- Create Sales Table
CREATE TABLE sales (
    id SERIAL PRIMARY KEY,
    movie_id INT REFERENCES movies(id) ON DELETE CASCADE,
    theater_id INT REFERENCES theaters(id) ON DELETE CASCADE,
    sale_date DATE NOT NULL,
    tickets_sold INT NOT NULL,
    ticket_price DECIMAL(5,2) NOT NULL
);

-- Insert Sample Theaters
INSERT INTO theaters (name) VALUES
('AMC Century City'),
('Regal LA Live'),
('Cinemark Playa Vista');

-- Insert Sample Movies
INSERT INTO movies (title) VALUES
('The Matrix'),
('Jurassic Park'),
('Titanic'),
('The Lion King'),
('Pulp Fiction'),
('Forrest Gump');

-- Insert Sample Sales Data
INSERT INTO sales (movie_id, theater_id, sale_date, tickets_sold, ticket_price)
VALUES
-- AMC Century City Sales
(1, 1, '2025-03-11', 50, 12.50), -- The Matrix
(1, 1, '2025-03-12', 30, 12.50),
(2, 1, '2025-03-11', 75, 15.00), -- Jurassic Park
(2, 1, '2025-03-12', 60, 15.00),

-- Regal LA Live Sales
(3, 2, '2025-03-11', 100, 10.00), -- Titanic
(3, 2, '2025-03-12', 90, 10.00),
(4, 2, '2025-03-11', 120, 8.50), -- The Lion King
(4, 2, '2025-03-12', 110, 8.50),

-- Cinemark Playa Vista Sales
(5, 3, '2025-03-11', 95, 14.00), -- Pulp Fiction
(5, 3, '2025-03-12', 85, 14.00),
(6, 3, '2025-03-11', 110, 13.00), -- Forrest Gump
(6, 3, '2025-03-12', 100, 13.00);

-- Verify inserted data
SELECT * FROM theaters;
SELECT * FROM movies;
SELECT * FROM sales;
