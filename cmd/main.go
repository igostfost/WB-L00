package main

import (
	"WB_L00/pkg/repository"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "4022",
		DBName:   "db1",
		SSLMode:  "disable",
	})

	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	//repos := repository.NewRepository(db)
	//service := services.NewService(repos)

	sc, err := connectNatsStream()
	defer sc.Close()

	// Публикация сообщения
	subject := "test-subject"
	message := "Messagegfds!!!FDS"

	subscription, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		fmt.Printf("Получено сообщение: %s\n", string(msg.Data))

		_, err = db.Exec("INSERT INTO messages (content) VALUES ($1)", string(msg.Data))
		if err != nil {
			log.Printf("Ошибка при записи в базу: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Ошибка при подписке на сообщения: %v", err)
	}
	defer subscription.Unsubscribe()
	for i := 0; i < 3; i++ {
		err = sc.Publish(subject, []byte(message))
		if err != nil {
			log.Fatalf("Ошибка при публикации сообщения: %v", err)
		}
		fmt.Printf("Сообщение опубликовано %d: %s\n", i, message)
	}

	// Ожидание сигнала для завершения программы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	fmt.Println("Программа завершена")

}

func connectNatsStream() (stan.Conn, error) {
	// Подключение к серверу NATS Streaming
	clusterID := "test-cluster"
	clientID := "client-123"
	natsURL := "nats://localhost:4222" // адрес сервера NATS

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL), stan.ConnectWait(time.Minute))
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
		return nil, err
	}
	return sc, nil
}
