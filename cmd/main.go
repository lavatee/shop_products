package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/endpoint"
	"github.com/lavatee/shop_products/internal/repository"
	"github.com/lavatee/shop_products/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"google.golang.org/grpc"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	if err := InitConfig(); err != nil {
		logrus.Fatalf("config error: %s", err.Error())
	}
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("env error: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.username"), os.Getenv("DB_PASSWORD"), viper.GetString("db.dbname"), viper.GetString("db.sslmode"))
	if err != nil {
		logrus.Fatalf("db error: %s", err.Error())
	}
	repo := repository.NewRepository(db, fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")))
	mqChannel, err := service.ConnectRabbitMQ(viper.GetString("rabbitmq.host"), viper.GetString("rabbitmq.port"), viper.GetString("rabbitmq.user"), os.Getenv("RABBIT_PASSWORD"))
	if err != nil {
		logrus.Fatalf("connect rabbitmq error: %s", err.Error())
	}
	producer := service.NewRabbitMQProducer(mqChannel)
	consumer := endpoint.NewRabbitMQConsumer(mqChannel)
	services := service.NewService(repo, producer)
	end := endpoint.NewEndpoint(services, consumer)
	srv := &products.Server{
		GRPCServer: grpc.NewServer(),
	}
	if err := end.Consumer.ConsumeQueue(service.PostDeleteProductEventQueue, end.ConsumePostDeleteProductEvent); err != nil {
		logrus.Fatalf("consume %s error: %s", service.PostDeleteProductEventQueue, err.Error())
	}
	if err := end.Consumer.ConsumeQueue(service.CompensateDeleteProductEventQueue, end.ConsumeCompensateDeleteProductEvent); err != nil {
		logrus.Fatalf("consume %s error: %s", service.CompensateDeleteProductEventQueue, err.Error())
	}
	if err := end.Consumer.ConsumeQueue(service.ConfirmProductDeletingQueue, end.ConsumeConfirmProductDeletingEvent); err != nil {
		logrus.Fatalf("consume %s error: %s", service.ConfirmProductDeletingQueue, err.Error())
	}
	go func() {
		if err := srv.Run(viper.GetString("port"), end); err != nil {
			logrus.Fatalf("run error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	srv.Shutdown()
	if err := db.Close(); err != nil {
		log.Fatalf("close db error: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
