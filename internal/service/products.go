package service

import (
	"errors"

	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/repository"
)

var categories = map[string]bool{
	"clothes":     true,
	"shoes":       true,
	"electronics": true,
	"toys":        true,
}

type ProductsService struct {
	Repo *repository.Repository
}

func NewProductsService(repo *repository.Repository) *ProductsService {
	return &ProductsService{
		Repo: repo,
	}
}

func (s *ProductsService) PostProduct(name string, amount int, price int, category string, description string, userId int) (int, error) {
	producer := PostProductProducer{
		Repo:      s.Repo,
		Observers: []PostProductObserver{},
	}
	event := producer.PostProduct(name, amount, price, category, description, userId)
	if !event.IsOk {
		return 0, errors.New(event.ErrorText)
	}
	return event.ProductId, nil
}

type PostProductEvent struct {
	ProductId          int
	ProductName        string
	ProductAmount      int
	ProductPrice       int
	ProductCategory    string
	ProductDescription string
	ProductCreatorId   int
	IsOk               bool
	ErrorText          string
}

type PostProductObserver interface {
	Update(event *PostProductEvent)
}

type PostProductProducer struct {
	Repo      *repository.Repository
	Observers []PostProductObserver
}

func (p PostProductProducer) PostProduct(name string, amount int, price int, category string, description string, userId int) PostProductEvent {
	if _, ok := categories[category]; !ok {
		return PostProductEvent{
			IsOk:      false,
			ErrorText: "invalid category",
		}
	}
	id, err := p.Repo.Products.PostProduct(name, amount, price, category, description, userId)
	if err != nil {
		return PostProductEvent{
			IsOk:      false,
			ErrorText: err.Error(),
		}
	}
	event := PostProductEvent{ProductId: id, ProductName: name, ProductAmount: amount, ProductPrice: price, ProductCategory: category, ProductDescription: description, ProductCreatorId: userId}
	for _, observer := range p.Observers {
		observer.Update(&event)
		if !event.IsOk {
			return event
		}
	}
	return event
}

func (s *ProductsService) DeleteProduct(id int) error {
	return s.Repo.DeleteProduct(id)
}

func (s *ProductsService) GetProducts(category string) ([]products.Product, error) {
	if _, ok := categories[category]; !ok {
		return nil, errors.New("invalid category")
	}
	return s.Repo.Products.GetProducts(category)
}

func (s *ProductsService) GetUserProducts(userId int) ([]products.Product, error) {
	return s.Repo.Products.GetUserProducts(userId)
}

func (s *ProductsService) GetSavedProducts(ids []int) ([]products.Product, error) {
	return s.Repo.GetSavedProducts(ids)
}

func (s *ProductsService) GetOneProduct(id int) (products.Product, error) {
	return s.Repo.GetOneProduct(id)
}
