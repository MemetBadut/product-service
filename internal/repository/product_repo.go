package repository

import (
	"github.com/MemetBadut/product-service/internal/model"
	"gorm.io/gorm"
)

// ProductRepository mendefinisikan kontrak operasi database.
// Menggunakan interface agar mudah di-mock saat testing.
type ProductRepository interface {
	Create(product *model.Product) error
	FindAll() ([]model.Product, error)
	FindByID(id uint) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	// GORM: INSERT INTO products (...) VALUES (...)
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	// GORM: SELECT * FROM products WHERE deleted_at IS NULL
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	// GORM: SELECT * FROM products WHERE id = ? AND deleted_at IS NULL
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *model.Product) error {
	// GORM: UPDATE products SET ... WHERE id = ?
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	// GORM soft delete: UPDATE products SET deleted_at = NOW() WHERE id = ?
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
	// Update stok dengan aman menggunakan atomic operation
	return r.db.Model(&model.Product{}).Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
