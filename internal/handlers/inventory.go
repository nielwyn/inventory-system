package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/service"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"go.uber.org/zap"
)

// InventoryHandler handles inventory endpoints
type InventoryHandler struct {
	inventoryService service.InventoryService
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// CreateItem handles creating a new inventory item
func (h *InventoryHandler) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, validator.FormatValidationError(err))
		return
	}

	item, err := h.inventoryService.CreateItem(&req)
	if err != nil {
		logger.Error("Failed to create item", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Item created successfully", item)
}

// GetAllItems handles retrieving all inventory items
func (h *InventoryHandler) GetAllItems(c *gin.Context) {
	items, err := h.inventoryService.GetAllItems()
	if err != nil {
		logger.Error("Failed to retrieve items", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve items")
		return
	}

	response.Success(c, http.StatusOK, "Items retrieved successfully", items)
}

// GetItemByID handles retrieving a single inventory item by ID
func (h *InventoryHandler) GetItemByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	item, err := h.inventoryService.GetItemByID(uint(id))
	if err != nil {
		logger.Error("Failed to retrieve item", zap.Error(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Item retrieved successfully", item)
}

// UpdateItem handles updating an inventory item
func (h *InventoryHandler) UpdateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, validator.FormatValidationError(err))
		return
	}

	item, err := h.inventoryService.UpdateItem(uint(id), &req)
	if err != nil {
		logger.Error("Failed to update item", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Item updated successfully", item)
}

// DeleteItem handles deleting an inventory item
func (h *InventoryHandler) DeleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.inventoryService.DeleteItem(uint(id)); err != nil {
		logger.Error("Failed to delete item", zap.Error(err))
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Item deleted successfully", nil)
}
