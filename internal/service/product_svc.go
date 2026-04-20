package service

import (
	"errors"

	"github.com/MemetBadut/product-service/internal/model"
	"github.com/MemetBadut/product-service/internal/repository"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(req *model.CreateProductRequest) (*model.Product, error)
	GetAllProducts() ([]model.Product, error)
	GetProductByID(id uint) (*model.Product, error)
	UpdateProduct(id uint, req *model.CreateProductRequest) (*model.Product, error)
	DeleteProduct(id uint) error
	CheckAndUpdateStock(id uint, quantity int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, errors.New("gagal membuat produk")
	}
	return product, nil
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetProductByID(id uint) (*model.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}
	return product, nil
}

func (s *productService) UpdateProduct(id uint, req *model.CreateProductRequest) (*model.Product, error) {
	product, err := s.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	if err := s.repo.Update(product); err != nil {
		return nil, errors.New("gagal update produk")
	}
	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	_, err := s.GetProductByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *productService) CheckAndUpdateStock(id uint, quantity int) error {
	product, err := s.GetProductByID(id)
	if err != nil {
		return err
	}
	if product.Stock < quantity {
		return errors.New("stok tidak mencukupi")
	}
	return s.repo.UpdateStock(id, -quantity)
}
