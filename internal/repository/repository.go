package repository

import (
	"github.com/jmoiron/sqlx"
	products "github.com/lavatee/shop_products"
)

type Repository struct {
	Products
}

type Products interface {
	PostProduct(name string, amount int, price int, category string, description string, userId int) (int, error)
	DeleteProduct(id int) error
	GetProducts(category string) ([]products.Product, error)
	GetUserProducts(userId int) ([]products.Product, error)
	GetSavedProducts(ids []int) ([]products.Product, error)
	GetOneProduct(id int) (products.Product, error)
	PostOrder(productId int) error
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Products: NewProductsPostgres(db),
	}
}
