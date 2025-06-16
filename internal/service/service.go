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
	PostDeleteProductEvent(eventId string, productId int, productCreator int) error
	GetProducts(category string) ([]products.Product, error)
	GetUserProducts(userId int) ([]products.Product, error)
	GetSavedProducts(ids []int) ([]products.Product, error)
	GetOneProduct(id int) (products.Product, error)
	CompensateDeleteProductEvent(compensatingEventId string, compensatedEventId string) error
	ConfirmProductDeleting(eventId string, confirmedEventId string) error
}

type MQProducer interface {
	SendMessage(queue string, message interface{}) error
	GetConfirmersAmount(queue string) (int, error)
}

func NewService(repo *repository.Repository, producer MQProducer) *Service {
	return &Service{
		Products: NewProductsService(repo, producer),
	}
}
