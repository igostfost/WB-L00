package main

import (
	"WB_L00"
	"WB_L00/pkg/handler"
	"WB_L00/pkg/repository"
	"WB_L00/pkg/services"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "4022",
		DBName:   "orderswb",
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	repos := repository.NewRepository(db)
	service := services.NewService(repos)

	sc, err := connectNatsStream()
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}
	defer sc.Close()

	cache := WB_L00.NewCache()

	ordersFromCache, err := service.GetAllOrdersFromDB()
	log.Printf("connect to cache and load data from db")

	fmt.Printf("order in cache: %v\n", len(ordersFromCache))
	//for _, v := range ordersFromCache {
	//	fmt.Println(v)
	//}
	handlers := handler.NewHandler(service, sc, cache)

	go func() {
		err = handlers.SubToChannel("orders")
		if err != nil {
			log.Fatalf("error subscribe to channel: %s", err)
		}
	}()

	serv := new(WB_L00.Server)

	go func() {
		err = serv.Run(":8000", handlers.InitRoutes())
		if err != nil {
			log.Fatalf("error start server: %s", err)
		}
	}()
	log.Println("Server started in port: 8000")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	if err := serv.Stop(context.Background()); err != nil {
		log.Fatalf("error stop server: %s", err)
	}
}

func connectNatsStream() (stan.Conn, error) {
	// Подключение к серверу NATS Streaming
	clusterID := "test-cluster"
	clientID := "client-123"
	natsURL := "nats://localhost:4222" // адрес сервера NATS

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.ConnectWait(time.Minute))
	if err != nil {
		return nil, err
	}
	return sc, nil
}
