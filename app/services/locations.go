package services

import (
	"receiptTracker/database"
)

type LocationServices struct {
	DB *database.DB
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *LocationServices) GetProductsByLocation(id int) ([]Product, error) {
	rows, err := s.DB.Query("
		SELECT p.id, p.name, p.category, p.barcode 
		FROM products p
		WHERE i.location_id = $1",
		id
	)
	if err != nil {
		return nil, err 
	}
	defer rows.Close()
	
	products := make([]Product, 0)

	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Category, &prod.Barcode); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}
	return products, nil
}

func (s *LocationServices) DeleteLocations() error {
	_, err := s.DB.Exec("DELETE FROM locations")
	return err
}

func (s *LocationServices) DeleteLocation(id int) error {
	return s.DB.DeleteByID("DELETE FROM locations WHERE id = $1", id)
}

func (s *LocationServices) GetLocations() ([]Location, error) {
	// Cursor pointing to results (doesnt alloc mem)
	// Effectively just a cursor, also a stream connection--whatever that means. 
	rows, err := s.DB.Query("SELECT id, name FROM locations")
	if err != nil { return nil, err }
	defer rows.Close()

	// Creates and empty slice, if there are no locations it
	// return []. var locations []Location would retunn null.
	locations := make([]Location, 0)

	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func (s *LocationServices) GetLocation(id int) (*Location, error) {
	var loc Location
	err := s.DB.GetByID(
		"SELECT id, name FROM locations WHERE id = $1",
		id,
		&loc.ID, &loc.Name, 
	)
	if err != nil { 
		return nil, err
	}
	return &loc, nil
}

func (s *LocationServices) CreateLocation(name string, locType string) (*Location, error) {
	id, err := s.DB.InsertReturningID(
		"INSERT INTO locations (name) VALUES ($1) RETURNING id",
		name,
	)
	if err != nil { 
		return nil, err 
	}
	return &Location{ID: id, Name: name}, nil
}

// ----> Examples/Refrences <----
// EXEC is for queries that dont return data (INSERT, DELETE, UPDATE)
// 		-> rowsAffected
//		-> LastInsertId
// Query for multiple rows
// QueryRow for a single row

