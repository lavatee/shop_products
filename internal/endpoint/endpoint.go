package endpoint

import (
	"github.com/lavatee/shop_products/internal/service"
	pb "github.com/lavatee/shop_protos/gen"
)

type Endpoint struct {
	pb.UnimplementedProductsServer
	Services *service.Service
}

func NewEndpoint(services *service.Service) *Endpoint {
	return &Endpoint{
		Services: services,
	}
}

// func (e *Endpoint) PostProduct(c context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {

// }
