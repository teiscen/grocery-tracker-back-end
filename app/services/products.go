package services

import (	
	"receiptTracker/database"
)

type ProductServices struct {
	DB *database.DB
}

// Should mirror the databasse 
type Product struct {
	ID 		 	int 	`json:"id"`
	Name 	 	string 	`json:"name"`
	LocationID 	*int 	`json:location_id`
}

// Delete/s
func (s *ProductServices) DeleteProducts() error {
	_, err := s.DB.Exec("DELETE FROM products")
	return err
}

func (s *ProductServices) DeleteProduct(id int) error {
	return s.DB.DeleteByID("DELETE FROM producs WHERE id = $1", id)
}

func (s *ProductServices) GetProducts() ([]Product, error) {
	products := make([]Product, 0)

	rows, err := s.DB.Query("SELECT id, name, location FROM products")
	if err != nil { 
		return nil, err 
	}
	defer rows.Close()

	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}
	return products, nil
}
 
// Get products by location is in locations
func (s *ProductServices) GetProduct(id int) (*Product, error) {
	var prod Product
	err := s.DB.GetByID(
		"SELECT id, name FROM products WHERE id = $1",
		id, 
		&prod.ID, &prod.Name,
	)
	if err != nil { 
		return nil, err 
	}
	return &prod, nil
}

func (s *ProductServices) EditProducts(id int) () {
	
}
 


func (s *ProductServices) CreateProduct(name string) (*Product, error) {
	id, err := s.DB.InsertReturningID(
		"INSERT INTO products (name, location) VALUES ($1) RETURNING id",
		name,
	)
	if err != nil { 
		return nil, err 
	}
	return &Product{ID: id, Name: name}, err
}