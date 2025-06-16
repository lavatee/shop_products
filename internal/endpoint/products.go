package endpoint

import (
	"context"
	"encoding/json"

	"github.com/lavatee/shop_products/internal/service"
	pb "github.com/lavatee/shop_protos/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (e *Endpoint) PostProduct(c context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	if req.Amount == 0 {
		return nil, status.Error(codes.InvalidArgument, "the request must have amount")
	}
	if req.Price == 0 {
		return nil, status.Error(codes.InvalidArgument, "the request must have price")
	}
	if req.Name == "" || req.Category == "" || req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "enter all fields")
	}
	productId, err := e.Services.Products.PostProduct(req.Name, int(req.Amount), int(req.Price), req.Category, req.Description, int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.PostProductResponse{
		Id: int64(productId),
	}, nil
}

func (e *Endpoint) GetProducts(c context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	if req.ProductCategory == "" {
		return nil, status.Error(codes.InvalidArgument, "the request must have category")
	}
	products, err := e.Services.Products.GetProducts(req.ProductCategory)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          int64(product.Id),
			UserId:      int64(product.UserId),
			Amount:      int64(product.Amount),
			Price:       int64(product.Price),
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
		}
	}
	return &pb.GetProductsResponse{
		Products: pbProducts,
	}, nil
}

func (e *Endpoint) GetUserProducts(c context.Context, req *pb.GetUserProductsRequest) (*pb.GetUserProductsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "the request must have user's id")
	}
	products, err := e.Services.Products.GetUserProducts(int(req.UserId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          int64(product.Id),
			UserId:      int64(product.UserId),
			Amount:      int64(product.Amount),
			Price:       int64(product.Price),
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
		}
	}
	return &pb.GetUserProductsResponse{
		Products: pbProducts,
	}, nil
}

func (e *Endpoint) GetOneProduct(c context.Context, req *pb.GetOneProductRequest) (*pb.GetOneProductResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "the request must have id")
	}
	product, err := e.Services.Products.GetOneProduct(int(req.Id))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetOneProductResponse{
		Product: &pb.Product{
			Id:          int64(product.Id),
			UserId:      int64(product.UserId),
			Amount:      int64(product.Amount),
			Price:       int64(product.Price),
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
		},
	}, nil
}

func (e *Endpoint) GetSavedProducts(c context.Context, req *pb.GetSavedProductsRequest) (*pb.GetSavedProductResponse, error) {
	if len(req.ProductsId) == 0 || req.ProductsId == nil {
		return nil, status.Error(codes.InvalidArgument, "the request must have ids")
	}
	ids := make([]int, len(req.ProductsId))
	for i, id := range req.ProductsId {
		ids[i] = int(id)
	}
	products, err := e.Services.Products.GetSavedProducts(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          int64(product.Id),
			UserId:      int64(product.UserId),
			Amount:      int64(product.Amount),
			Price:       int64(product.Price),
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
		}
	}
	return &pb.GetSavedProductResponse{
		Products: pbProducts,
	}, nil
}

func (e *Endpoint) ConsumePostDeleteProductEvent(data []byte) error {
	var event service.PostDeleteProductEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logrus.Errorf("Error in ConsumePostDeleteProductEvent: %s", err.Error())
		return err
	}
	if err := e.Services.Products.PostDeleteProductEvent(event.EventId, event.EventInfo.ProductId, event.EventInfo.ProductCreator); err != nil {
		logrus.Errorf("Error in ConsumePostDeleteProductEvent: %s", err.Error())
		return err
	}
	logrus.Infof("Event with id '%s' was consumed successfully!", event.EventId)
	return nil
}

func (e *Endpoint) ConsumeCompensateDeleteProductEvent(data []byte) error {
	var event service.CompensateDeleteProductEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logrus.Errorf("Error in ConsumeCompensateDeleteProductEvent: %s", err.Error())
		return err
	}
	if err := e.Services.Products.CompensateDeleteProductEvent(event.EventId, event.EventInfo.DeleteProductEventId); err != nil {
		logrus.Errorf("Error in ConsumeCompensateDeleteProductEvent: %s", err.Error())
		return err
	}
	logrus.Infof("Event with id '%s' was consumed successfully!", event.EventId)
	return nil
}

func (e *Endpoint) ConsumeConfirmProductDeletingEvent(data []byte) error {
	var event service.ConfirmProductDeletingEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logrus.Errorf("Error in ConsumeConfirmProductDeletingEvent: %s", err.Error())
		return err
	}
	if err := e.Services.Products.ConfirmProductDeleting(event.EventId, event.EventInfo.DeleteProductEventId); err != nil {
		logrus.Errorf("Error in ConsumeConfirmProductDeletingEvent: %s", err.Error())
		return err
	}
	logrus.Info("Event with id '%s' was consumed successfully!", event.EventId)
	return nil
}

func (e *Endpoint) ConsumePostPostOrderEvent(data []byte) error {

}
