package endpoint

import (
	"context"

	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/service"
	pb "github.com/lavatee/shop_protos/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (e *Endpoint) PostProduct(c context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	if req.Amount == 0 {
		return nil, status.Error(codes.InvalidArgument, "amount can't be 0")
	}
	if req.Price == 0 {
		return nil, status.Error(codes.InvalidArgument, "price can't be 0")
	}
	if req.Name == "" || req.Category == "" || req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "enter all fields")
	}
	producer := PostProductProducer{
		Services:  e.Services,
		Observers: []PostProductObserver{},
	}
	event := producer.Produce(req.Name, int(req.Amount), int(req.Price), req.Category, req.Description, int(req.UserId))
	if !event.IsOk {
		return nil, status.Error(codes.Internal, event.Error)
	}
	return &pb.PostProductResponse{
		Id: int64(event.Id),
	}, nil
}

type PostProductEvent struct {
	Name        string
	Amount      int
	Price       int
	Category    string
	Description string
	Id          int
	UserId      int
	IsOk        bool
	Error       string
}

type PostProductProducer struct {
	Observers []PostProductObserver
	Services  *service.Service
}

type PostProductObserver interface {
	Update(event *PostProductEvent)
}

func (p PostProductProducer) Produce(name string, amount int, price int, category string, description string, userId int) PostProductEvent {
	id, err := p.Services.PostProduct(name, amount, price, category, description, userId)
	if err != nil {
		return PostProductEvent{IsOk: false, Error: err.Error()}
	}
	event := PostProductEvent{
		Name:        name,
		Amount:      amount,
		Price:       price,
		Category:    category,
		Description: description,
		UserId:      userId,
		Id:          id,
		IsOk:        true,
	}
	for _, observer := range p.Observers {
		observer.Update(&event)
		if !event.IsOk {
			return event
		}
	}
	return event
}

func (e *Endpoint) GetProducts(c context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	if req.ProductCategory == "" {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}
	producer := GetProductsProducer{
		Services:  e.Services,
		Observers: []GetProductObserver{},
	}
	event := producer.Produce(req.ProductCategory)
	if !event.IsOk {
		return nil, status.Error(codes.Internal, event.Error)
	}
	return &pb.GetProductsResponse{
		Products: event.GRPCProducts,
	}, nil
}

type GetProductsEvent struct {
	Category     string
	Products     []products.Product
	GRPCProducts []*pb.Product
	IsOk         bool
	Error        string
}

type GetProductsProducer struct {
	Observers []GetProductObserver
	Services  *service.Service
}

type GetProductObserver interface {
	Update(event *GetProductsEvent)
}

func (p GetProductsProducer) Produce(category string) GetProductsEvent {
	products, err := p.Services.GetProducts(category)
	if err != nil {
		return GetProductsEvent{IsOk: false, Error: err.Error()}
	}
	gRPCProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		gRPCProducts[i] = &pb.Product{
			Id:          int64(product.Id),
			UserId:      int64(product.UserId),
			Price:       int64(product.Price),
			Amount:      int64(product.Amount),
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
		}
	}
	event := GetProductsEvent{
		Products:     products,
		GRPCProducts: gRPCProducts,
		Category:     category,
		IsOk:         true,
	}
	for _, observer := range p.Observers {
		observer.Update(&event)
		if !event.IsOk {
			return event
		}
	}
	return event
}
