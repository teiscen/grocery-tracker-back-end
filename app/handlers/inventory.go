package handlers

import(
	"net/http"
	"receiptTracker/services"
)

type InventoryHandler struct {
	Service *services.InventoryServices
}

func (h *InventoryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method{
		case http.MethodGet:
			h.GetInventory(w, r)
		case http.MethodPost:
			h.CreateInventoryItem(w, r)
		default:
			http.Error(w, "Method nto allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/inventory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method{
		case http.MethodGet:
			h.GetInventoryItem(w, r)
		case http.MethodPost:
			h.UpdateInventoryItem(w, r)
		case http.MethodDelete:
			h.DeleteInventoryItem(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// GET /api/inventory/3
func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/api/inventory/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory id", err)
		return
	}
	item, err := h.Service.GetInventoryItem(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get inventory item", err)
		return
	}
	writeJson(w, http.StatusOK, item)
}

// POST /api/inventory
func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ProductID  int     `json:"productId"`
		LocationID int     `json:"locationId"`
		Quantity   float64 `json:"quantity"`
		Unit       string  `json:"unit"`
		ExpiryDate *string `json:"expiryDate"`
		Opened     bool    `json:"opened"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	item, err := h.Service.CreateInventoryItem(
		body.ProductID,
		body.LocationID,
		body.Quantity,
		body.Unit,
		body.ExpiryDate,
		body.Opened,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create inventory item", err)
		return
	}
	writeJson(w, http.StatusCreated, item)
}

// PUT /api/inventory/3
func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/api/inventory/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid inventory id", err)
		return
	}
	var body struct {
		LocationID int     `json:"locationId"`
		Quantity   float64 `json:"quantity"`
		Unit       string  `json:"unit"`
		ExpiryDate *string `json:"expiryDate"`
		Opened     bool    `json:"opened"`
	}
	if err := decodeBody(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	item, err := h.Service.UpdateInventoryItem(
		id,
		body.LocationID,
		body.Quantity,
		body.Unit,
		body.ExpiryDate,
		body.Opened,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update inventory item", err)
		return
	}
	writeJson(w, http.StatusOK, item)
}

// DELETE /api/inventory/3
func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/api/inventory/")
	if err != nil {
		writeerror(w, http.statusbadrequest, "invalid inventory id", err)
		return
	}
	err = h.Service.DeleteInventoryItem(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete inventory item", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
