package service

import (
	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/repository"
)

type Service struct {
	Products
}

type Products interface {
	PostProduct(name string, amount int, price int, category string, description string, userId int) (int, error)
	DeleteProduct(id int) error
	GetProducts(category string) ([]products.Product, error)
	GetUserProducts(userId int) ([]products.Product, error)
	GetSavedProducts(ids []int) ([]products.Product, error)
	GetOneProduct(id int) (products.Product, error)
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Products: NewProductsService(repo),
	}
}
