

CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Represents a product ttype that exists in hte world (Dairyland Whole Milk 1L)
-- Barcode is unique, two products cant share one 
-- TODO: warn users that changes affect all instances of the profuct 
CREATE TABLE product (
    id SERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL, 
    category VARCHAR(50), 
    barcode VARCHAR(50) UNIQUE
);

-- Represents the actual instance of a product inside of a location
-- Many to one relationship with product (1 product may have many inventories)
-- Cannot delete a location or product if inventory refrences the
-- TODO: Decide on cascading deletes
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY, 
    product_id INT NOT NULL REFERENCES products(id),
    location_id INT NOT NULL REFERENCES locations(id),
    quantity DECIMAL(10,2) NOT NULL, 
    unit VARCHAR(20) NOT NULL,
    expiry_date DATE,
    opened BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

