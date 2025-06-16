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
	PostDeleteProductEvent(id int, eventId string, productCreator int) error
	GetProducts(category string) ([]products.Product, error)
	GetUserProducts(userId int) ([]products.Product, error)
	GetSavedProducts(ids []int) ([]products.Product, error)
	GetOneProduct(id int) (products.Product, error)
	PostEvent(eventId string, eventType string, confirmersAmount int) error
	PostCompensatoryEvent(comEventId string, eventId string) error
	CompensateDeleteProductEvent(eventId string) error
	ConfirmProductDeleting(eventId string) error
}

func NewRepository(db *sqlx.DB, redisHost string) *Repository {
	return &Repository{
		Products: NewProductsPostgres(db, NewRedisDB(redisHost)),
	}
}
