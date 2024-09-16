package api

import (
	"encoding/json"
	"net/http"

	"github.com/Drynok/tx-parser/internal/parser"
)

type Handler struct {
	parser parser.Parser
}

func NewHandler(p parser.Parser) *Handler {
	return &Handler{parser: p}
}

func (h *Handler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := h.parser.GetCurrentBlock()
	json.NewEncoder(w).Encode(map[string]int{"current_block": block})
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	success := h.parser.Subscribe(req.Address)
	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}
	transactions := h.parser.GetTransactions(address)
	json.NewEncoder(w).Encode(transactions)
}
