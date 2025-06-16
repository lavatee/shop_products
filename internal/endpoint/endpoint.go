package endpoint

import (
	"github.com/lavatee/shop_products/internal/service"
	pb "github.com/lavatee/shop_protos/gen"
)

type MQConsumer interface {
	ConsumeQueue(queue string, handler func([]byte) error) error
}

type Endpoint struct {
	pb.UnimplementedProductsServer
	Services *service.Service
	Consumer MQConsumer
}

func NewEndpoint(services *service.Service, consumer MQConsumer) *Endpoint {
	return &Endpoint{
		Services: services,
		Consumer: consumer,
	}
}

// func (e *Endpoint) PostProduct(c context.Context, req *pb.PostProductRequest) (*pb.PostProductResponse, error) {

// }
