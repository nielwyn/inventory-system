package service

import (
	"errors"
	"fmt"

	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/repository"
)

// InventoryService handles inventory business logic
type InventoryService interface {
	CreateItem(req *models.CreateItemRequest) (*models.Item, error)
	GetAllItems() ([]models.Item, error)
	GetItemByID(id uint) (*models.Item, error)
	UpdateItem(id uint, req *models.UpdateItemRequest) (*models.Item, error)
	DeleteItem(id uint) error
}

type inventoryService struct {
	repo repository.InventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(repo repository.InventoryRepository) InventoryService {
	return &inventoryService{repo: repo}
}

// CreateItem creates a new inventory item
func (s *inventoryService) CreateItem(req *models.CreateItemRequest) (*models.Item, error) {
	// Check if SKU already exists
	existingItem, err := s.repo.FindBySKU(req.SKU)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		return nil, errors.New("item with this SKU already exists")
	}

	// Create item
	item := &models.Item{
		Name:        req.Name,
		SKU:         req.SKU,
		Description: req.Description,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Category:    req.Category,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetAllItems retrieves all inventory items
func (s *inventoryService) GetAllItems() ([]models.Item, error) {
	return s.repo.FindAll()
}

// GetItemByID retrieves an item by ID
func (s *inventoryService) GetItemByID(id uint) (*models.Item, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}
	return item, nil
}

// UpdateItem updates an existing item
func (s *inventoryService) UpdateItem(id uint, req *models.UpdateItemRequest) (*models.Item, error) {
	// Find existing item
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}

	// Check if SKU is being updated and if it already exists
	if req.SKU != nil && *req.SKU != item.SKU {
		existingItem, err := s.repo.FindBySKU(*req.SKU)
		if err != nil {
			return nil, err
		}
		if existingItem != nil {
			return nil, fmt.Errorf("item with SKU '%s' already exists", *req.SKU)
		}
		item.SKU = *req.SKU
	}

	// Update fields if provided
	if req.Name != nil {
		item.Name = *req.Name
	}
	if req.Description != nil {
		item.Description = *req.Description
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Category != nil {
		item.Category = *req.Category
	}

	// Save updated item
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem deletes an item by ID
func (s *inventoryService) DeleteItem(id uint) error {
	// Check if item exists
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}

	return s.repo.Delete(id)
}
