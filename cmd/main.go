package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/endpoint"
	"github.com/lavatee/shop_products/internal/repository"
	"github.com/lavatee/shop_products/internal/service"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("config error", err.Error())
	}
	if err := gotenv.Load(); err != nil {
		log.Fatalf("env error", err.Error())
	}
	db, err := repository.NewPostgresDB(viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.username"), os.Getenv("DB_PASSWORD"), viper.GetString("db.dbname"), viper.GetString("db.sslmode"))
	if err != nil {
		log.Fatalf("db error", err.Error())
	}
	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	end := endpoint.NewEndpoint(services)
	srv := &products.Server{
		GRPCServer: grpc.NewServer(),
	}
	go func() {
		if err := srv.Run(viper.GetString("port"), end); err != nil {
			log.Fatalf("run error", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	srv.Shutdown()
	if err := db.Close(); err != nil {
		log.Fatalf("close db error", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
