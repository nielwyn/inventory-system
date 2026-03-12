package repository

import (
	"errors"

	"github.com/nielwyn/inventory-system/internal/models"
	"gorm.io/gorm"
)

// InventoryRepository handles inventory data operations
type InventoryRepository interface {
	Create(item *models.Item) error
	FindAll() ([]models.Item, error)
	FindByID(id uint) (*models.Item, error)
	FindBySKU(sku string) (*models.Item, error)
	Update(item *models.Item) error
	Delete(id uint) error
}

type inventoryRepository struct {
	db *gorm.DB
}

// NewInventoryRepository creates a new inventory repository
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Create creates a new item
func (r *inventoryRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

// FindAll retrieves all items
func (r *inventoryRepository) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}

// FindByID finds an item by ID
func (r *inventoryRepository) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// FindBySKU finds an item by SKU
func (r *inventoryRepository) FindBySKU(sku string) (*models.Item, error) {
	var item models.Item
	err := r.db.Where("sku = ?", sku).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Update updates an existing item
func (r *inventoryRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

// Delete soft deletes an item by ID
func (r *inventoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Item{}, id).Error
}
