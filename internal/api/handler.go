package api

import (
	"net/http"

	"github.com/Drynok/tx-parser/internal/parser"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	parser parser.Parser
}

func NewHandler(p parser.Parser) *Handler {
	return &Handler{parser: p}
}

func (h *Handler) GetCurrentBlock(c *gin.Context) {
	block := h.parser.GetCurrentBlock()
	c.JSON(http.StatusOK, gin.H{"current_block": block})
}

func (h *Handler) Subscribe(c *gin.Context) {
	var req struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	success := h.parser.Subscribe(req.Address)
	c.JSON(http.StatusOK, gin.H{"success": success})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address is required"})
		return
	}
	transactions := h.parser.GetTransactions(address)
	c.JSON(http.StatusOK, transactions)
}
