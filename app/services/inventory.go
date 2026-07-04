package services

import(
	"receiptTracker/database"
	"time"
)

type InventoryServices struct {
	DB *database.DB
}

// computed fields (expiryStatus, daysUntilExpiry) are calculated before returning
type InventoryItem struct {
    ID             int     `json:"id"`
    ProductID      int     `json:"productId"`
    Name           string  `json:"name"`
    Category       *string `json:"category,omitempty"`
    Barcode        *string `json:"barcode,omitempty"`
    LocationID     int     `json:"locationId"`
    LocationName   string  `json:"locationName"`
    Quantity       float64 `json:"quantity"`
    Unit           string  `json:"unit"`
    ExpiryDate     *string `json:"expiryDate,omitempty"`
    DaysUntilExpiry *int   `json:"daysUntilExpiry,omitempty"`
    Opened         bool    `json:"opened"`
    ExpiryStatus   string  `json:"expiryStatus"`
}

// Defaults to ("ok", nil)
func computeExpiry(expiryDate *string) (string, *int) {
	if expiryDate == nil {
		return "ok", nil
	}
	t, err = time.Parse(time.DateOnly, *expiryDate)
	if err != nil {
		return "ok", nil
	}
	
	days := int(time.Until(t).Hours() / 24) 	
	switch {
	case days < 0:
		return "expired", &days
	case days <=3:
		return "expiring", &days
	default:
		return "ok", &days
	}
}

const baseInventoryQuery = `
	SELECT
		i.id,
		i.quantity,
		i.unit,
		i.opened,
		i.expiry_date,
		p.id   AS product_id,
		p.name,
		p.category,
		p.barcode,
		l.id   AS location_id,
		l.name AS location_name
	FROM inventory i
	JOIN products  p ON i.product_id  = p.id
	JOIN locations l ON i.location_id = l.id`


func (s *InventoryServices) GetInventory(LocationID *int) ([]InventoryItem, error) {
	query := baseInventoryQuery 

	// Filter by location if its provided
	args := make([]any, 0)
	if locationID != nil {
		query += " WHERE l.id = $1"
		args = append(args, *locationID)
	}
	
	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]InventoryItem, 0)
	for rows.Next() {
		var item InventoryItem
		var expiryDate *string 

		if err := rows.Scan(
			&item.ID,
			&item.Quantity, 
			&item.Unit,
			&item.Opened,
			&expiryDate,
			&item.ProductID,
			&item.Name, 
			&item.Category, 
			&item.Barcode, 
			&item.LocationID,
			&item.LocationName,
		); err != nil {
			return nil, err
		}

		item.ExpiryDate = expiryDate
		item.ExpiryStatus, item.DaysUntilExpiry = computeExpiry(expiryDate)
		items = append(items, item)
	}
	return items, nil
}

func (s *InventoryServices) GetInventoryItem(id int) (*InventoryItem, error) {
	var item InventoryItem
	var expiryDate *string

	query := (baseInventoryQuery += " Where i.id = $1")
	err := s.DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.Quantity, 
		&item.Unit,
		&item.Opened,
		&expiryDate,
		&item.ProductID,
		&item.Name, 
		&item.Category, 
		&item.Barcode, 
		&item.LocationID,
		&item.LocationName,
	)
	if err != nil {
		return nil, err
	}

	item.ExpiryDate = expiryDate
	item.ExpiryStatus, item.DaysUntilExpiry = computeExpiry(expiryDate)
	return &item, nil
}

func (s *InventoryServices) CreateInventoryItem(
	name string, category *string, barcode *string, locationID int,
	quantity float64, unit string, expiryDate *string, opened bool,
) (*InventoryItem, error) {
	id, err := s.DB.InsertReturningId(`
		INSERT INTO inventory (product_id, location_id, quantity, unit, expiry_date,
		VALUES $1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		productID, locationID, quantity, unit, expiryDate, opened,
	)
	if err != nil {
		return nil, err
	}
	return s.GetInventoryItem(id)
}

