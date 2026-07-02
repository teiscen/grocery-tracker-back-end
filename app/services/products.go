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
	Category	string 	`json:"category"`
	Barcode		string 	`json:"barcode"`
}

// Delete/s
func (s *ProductServices) DeleteProducts() error {
	_, err := s.DB.Exec("DELETE FROM products")
	return err
}

func (s *ProductServices) DeleteProduct(id int) error {
	return s.DB.DeleteByID("DELETE FROM products WHERE id = $1", id)
}

func (s *ProductServices) GetProducts() ([]Product, error) {
	products := make([]Product, 0)

	rows, err := s.DB.Query("SELECT id, name, category, barcode location FROM products")
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
		"SELECT id, name, category, barcode FROM products WHERE id = $1",
		id, 
		&prod.ID, &prod.Name,
	)
	if err != nil { 
		return nil, err 
	}
	return &prod, nil
}

// Returns product if found, or nil, nil if it was not found by barcode
func (s *ProductServices) GetProductByBarcode(barcode string) (*Product, error) {
	var prod Product
	err := s.DB.QueryRow(
		"SELECT id, name, category, barcode FROM products WHERE barcode = $1",
		barcode,
	).Scan(&prod.ID, &prod.Name, &prod.Category, &prod.Barcode)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	} 
	return &prod, nil
}

func (s *ProductServices) SearchProducts(query string) ([]Product, error) {
	rows, err := s.DB.Query( 										// Case insentive like
		"SELECT id, name, category, barcode FROM products WHERE name ILIKE $1",
		"%"+query+"%", // % Wildcard match anything thing 
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Category, &prod.Barcode); err != nil{
			return nil, err
		}
		products = append(products, prod)
	}
	return products, nil
}

// func (s *ProductServices) EditProducts(id int) () {
// }

func (s *ProductServices) CreateProduct(name string, category string, barcode string) (*Product, error) {
	id, err := s.DB.InsertReturningID(
		"INSERT INTO products (name, category, barcode) VALUES ($1, $2, $3) RETURNING id",
		name, category, barcode, 
	)
	if err != nil { 
		return nil, err 
	}
	return &Product{ID: id, Name: name, Category category, Barcode barcode}, nil
}
