-- Products
CREATE TABLE products (
    id       SERIAL PRIMARY KEY, 
    name     VARCHAR(100) NOT NULL, 
    location_id INT REFERENCES locations(id)
)
-- Locations
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

----> Example/Reference: <----

-- products
-- CREATE TABLE products (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     category VARCHAR(50),
--     brand VARCHAR(100),
--     barcode VARCHAR(50),
--     unit VARCHAR(20) NOT NULL,
--     default_location_id INT REFERENCES locations(id),
--     restock_threshold INT DEFAULT 1,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- inventory
-- CREATE TABLE inventory (
--     id SERIAL PRIMARY KEY,
--     product_id INT NOT NULL REFERENCES products(id),
--     location_id INT NOT NULL REFERENCES locations(id),
--     quantity DECIMAL(10,2) NOT NULL,
--     expiry_date DATE,
--     opened BOOLEAN DEFAULT FALSE,
--     date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     notes TEXT
-- );
