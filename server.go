package products

import (
	"fmt"
	"net"

	pb "github.com/lavatee/shop_protos/gen"
	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
}

func (s *Server) Run(port string, handler pb.ProductsServer) error {
	pb.RegisterProductsServer(s.GRPCServer, handler)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	if err := s.GRPCServer.Serve(listener); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown() {
	s.GRPCServer.GracefulStop()
}
