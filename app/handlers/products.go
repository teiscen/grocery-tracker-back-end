package handlers

import (
	"net/http"

	"receiptTracker/services"
)

type ProductHandler struct {
	Service *services.ProductServices
	// createProduct
	// getProduct/s
	// deleteProduct/s
}

func (h *ProductHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:	h.GetProducts(w, r)
		case http.MethodPost:	h.CreateProduct(w, r)
		case http.MethodDelete: h.DeleteProducts(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}		
	})

	mux.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet: 	h.GetProduct(w, r)
		case http.MethodDelete: h.DeleteProduct(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}		
	})
}

func (h *ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	err := h.Service.DeleteProducts()
	if err != nil {
		writeError(w, http.StatusInternalServerError,
				   	"delete produts error", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/product/")
	if err != nil {
		writeError(w, http.StatusBadRequest,
				   	"invalid product id", err)
		return
	}

	err = h.Service.DeleteProduct(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
				  "failed to delete product", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}	

func (h *ProductHandler) getProductByBarcode(w http.ResponseWriter, r *http.Request, string barcode) {
	product, err := h.Service.GetProductByBarcode(barcode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "barcode lookup failed", err)
		return 
	}
	if product == nil {
		// Not found
		writeJson(w, http.StatusNotFound, map[string]string{"message": "product not found"})
		return
	}
	writeJson(w, http.StatusOK, product)
	return
}

func (h *ProductHandler) getProductBySearch(search string) () {
	products, err := h.Service.SearchProducts(search)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "search failed", err)
		return 
	}	
	writeJson(w, http.StatusOK, products)
	return
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	barcode := r.URL.Query().Get("barcode")
	search := r.URL.Query().get("search")

	if barcode != "" {
		h.getProductByBarcode(w, r, barcode)
		return
	}

	if search != "" {
		h.getProductBySearch(w, r, search)
		return 
	}

	products, err := h.Service.GetProducts()
	if err != nil { 
		writeError(w, http.StatusInternalServerError,
					"get products error", err)
		return
	}
	writeJson(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path, "/product")
	if err != nil {
		writeError(w, http.StatusBadRequest, 
					"invalid product id", err)
		return 
	}

	product, err := h.Service.GetProduct(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError,
					"failed to get product", err)
		return
	}
	writeJson(w, http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// ToDo: Move out
	var body struct {
		Name string `json:"name"`
	}
	if err := decodeBody(r, &body); err != nil {
        writeError(w, http.StatusBadRequest, 
					"invalid request body", err)
        return
	}
	product, err := h.Service.CreateProduct(body.Name)
	if err != nil {
        writeError(w, http.StatusInternalServerError, 
					"failed to create product", err)
        return
	}
    writeJson(w, http.StatusCreated, product)
}
