package main

import (
	"fmt"
	"net/http"

	"receiptTracker/database"
	"receiptTracker/handlers"
	"receiptTracker/services"
)

func registerRoutes(resources []handlers.Registerable) *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range resources {
		r.Register(mux)
	}
	return mux
}

func main() {
	// Establish DB connection, and ensure it closes
	db := database.Connect()
	defer db.Close()

	// Initialize Services
	locationService := &services.LocationServices{
		DB: db,
	}
	productService  := &services.ProductServices{
		DB: db,
	}
	inventoryService := &services.InventoryServices{
		DB: db,
	}
	// Initialize Handlers with resource pattern
	mux := registerRoutes([]handlers.Registerable{
		&handlers.LocationHandler{	Service: locationService},
		&handlers.ProductHandler{	Service: productService},
		&handlers.InventoryHandler{ Service: inventoryService},
	})

	fmt.Println("starting server on port 8000...")
	http.ListenAndServe(":8000", mux)
}
