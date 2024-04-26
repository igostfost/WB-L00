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
		Username: "test_user",
		Password: "qwerty",
		DBName:   "orderswb",
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	repos := repository.NewRepository(db)
	service := services.NewService(repos)

	sc, err := connectNatsStream()
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}

	ordersFromCache, err := service.GetAllOrdersFromDB()
	if err != nil {
		log.Fatalf("Failed to get all orders from db: ", err)

	}

	cache := WB_L00.NewCache()
	cache.RestoreFromDB(ordersFromCache)
	log.Printf("connect to cache and get all orders from db")

	fmt.Printf("order in cache: %v\n", len(ordersFromCache))

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
			log.Printf("error start server: %s", err)
			return
		}
	}()

	log.Println("Server started in port: 8000")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	if err := serv.Stop(context.Background()); err != nil {
		log.Printf("error stop server: %s", err)
	}

	if err := sc.Close(); err != nil {
		log.Printf("error close nats streaming client: %s", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("error close db: %s", err)
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
