package service

import (
	"errors"
	"fmt"

	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/repository"
)

var categories = map[string]bool{
	"clothes":     true,
	"shoes":       true,
	"electronics": true,
	"toys":        true,
}

type DeleteProductEventInfo struct {
	ProductId      int `json:"productId"`
	ProductCreator int `json:"productCreator"`
}

type PostDeleteProductEvent struct {
	EventInfo DeleteProductEventInfo `json:"eventInfo"`
	EventId   string                 `json:"eventId"`
}

type CompensateDeleteProductEventInfo struct {
	ProductId            int    `json:"productId"`
	DeleteProductEventId string `json:"deleteProductEventId"`
}

type CompensateDeleteProductEvent struct {
	EventInfo CompensateDeleteProductEventInfo `json:"eventInfo"`
	EventId   string                           `json:"eventId"`
}

type ConfirmProductDeletingEventInfo struct {
	DeleteProductEventId string `json:"deleteProductEventId"`
}

type ConfirmProductDeletingEvent struct {
	EventInfo ConfirmProductDeletingEventInfo `json:"eventInfo"`
	EventId   string                          `json:"eventId"`
}

type PostOrderEventProduct struct {
	ProductId int `json:"productId"`
	Amount    int `json:"amount"`
	Price     int `json:"price"`
}

type PostOrderEventInfo struct {
	UserId   int                     `json:"userId"`
	Price    int                     `json:"price"`
	Products []PostOrderEventProduct `json:"products"`
}

type PostOrderEvent struct {
	EventInfo PostOrderEventInfo `json:"eventInfo"`
	EventId   string             `json:"eventId"`
}

type CompensatePostOrderEventInfo struct {
	OrderingEventId string `json:"orderingEventId"`
}

type CompensatePostOrderEvent struct {
	EventInfo CompensatePostOrderEventInfo `json:"eventInfo"`
	EventId   string                       `json:"eventId"`
}

type ConfirmPostOrderEventInfo struct {
	OrderingEventId string `json:"orderingEventId"`
}

type ConfirmPostOrderEvent struct {
	EventInfo ConfirmPostOrderEventInfo `json:"eventInfo"`
	EventId   string                    `json:"eventId"`
}

const (
	PostDeleteProductEventQueue       = "post_delete_product_event"
	CompensateDeleteProductEventQueue = "compensate_delete_product_event"
	DeleteProductType                 = "delete_product"
	ConfirmProductDeletingQueue       = "confirm_product_deleting"
	ConfirmEventType                  = "confirm_event"
)

type ProductsService struct {
	Repo     *repository.Repository
	Producer MQProducer
}

func NewProductsService(repo *repository.Repository, producer MQProducer) *ProductsService {
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
	event := PostProductEvent{ProductId: id, ProductName: name, ProductAmount: amount, ProductPrice: price, ProductCategory: category, ProductDescription: description, ProductCreatorId: userId, IsOk: true}
	for _, observer := range p.Observers {
		observer.Update(&event)
		if !event.IsOk {
			return event
		}
	}
	return event
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

func (s *ProductsService) PostDeleteProductEvent(eventId string, productId int, productCreator int) error {
	compensatoryEvent := CompensateDeleteProductEventInfo{
		ProductId:            productId,
		DeleteProductEventId: eventId,
	}
	confirmersAmount, err := s.Producer.GetConfirmersAmount(PostDeleteProductEventQueue)
	if err != nil {
		if err := s.Producer.SendMessage(CompensateDeleteProductEventQueue, compensatoryEvent); err != nil {
			return fmt.Errorf("Get confirmers error: Send compensatory delete product event error: %s", err.Error())
		}
		return fmt.Errorf("Get confirmers error (compensatory event was sent successfully): %s", err.Error())
	}
	if err := s.Repo.Products.PostEvent(eventId, DeleteProductType, confirmersAmount); err != nil {
		if err := s.Producer.SendMessage(CompensateDeleteProductEventQueue, compensatoryEvent); err != nil {
			return fmt.Errorf("Post event error: Send compensatory delete product event error: %s", err.Error())
		}
		return fmt.Errorf("Post event error (compensatory event was sent successfully): %s", err.Error())
	}
	if err := s.Repo.Products.PostDeleteProductEvent(productId, eventId, productCreator); err != nil {
		if err := s.Producer.SendMessage(CompensateDeleteProductEventQueue, compensatoryEvent); err != nil {
			return fmt.Errorf("Post delete product event error: Send compensatory delete product event error: %s", err.Error())
		}
		return fmt.Errorf("Post delete product event error (compensatory event was sent successfully): %s", err.Error())
	}
	message := ConfirmProductDeletingEventInfo{
		DeleteProductEventId: eventId,
	}
	if err := s.Producer.SendMessage(ConfirmProductDeletingQueue, message); err != nil {
		return fmt.Errorf("Post delete product event error: Send confirm product deleting event error: %s", err.Error())
	}
	return nil
}

func (s *ProductsService) CompensateDeleteProductEvent(compensatingEventId string, compensatedEventId string) error {
	if err := s.Repo.Products.PostCompensatoryEvent(compensatingEventId, compensatedEventId); err != nil {
		return fmt.Errorf("Compensate delete product event error (post compensatory event): %s", err.Error())
	}
	if err := s.Repo.Products.CompensateDeleteProductEvent(compensatedEventId); err != nil {
		return fmt.Errorf("Compensate delete product event error (change product's status): %s", err.Error())
	}
	return nil
}

func (s *ProductsService) ConfirmProductDeleting(eventId string, confirmedEventId string) error {
	if err := s.Repo.Products.PostEvent(eventId, ConfirmEventType, 0); err != nil {
		return fmt.Errorf("Event with id '%s' has already been consumed", eventId)
	}
	if err := s.Repo.Products.ConfirmProductDeleting(confirmedEventId); err != nil {
		return fmt.Errorf("Confirm product deleting error: %s", err.Error())
	}
	return nil
}

//Сделать PostDeleteProductEvent и CompensateDeleteProductEvent с паттерном Observer
