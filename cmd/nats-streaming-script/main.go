package main

import (
	"WB_L00/types"
	"encoding/json"
	"github.com/icrowley/fake"
	"math/rand"
	"time"

	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	// Подключение к серверу NATS Streaming
	clusterID := "test-cluster"
	clientID := "client_2"
	natsURL := "nats://localhost:4222" // адрес сервера NATS

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// генерим заказы и отправляем в канал, принтуем UID каждого заказа
	for i := 0; i < 50; i++ {
		order := createOrders()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		channel := "orders"
		err = sc.Publish(channel, orderJSON)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Order data sent to channel")
}

func createOrders() types.Order {

	order := types.Order{
		OrderUID:          fake.DigitsN(13),
		TrackNumber:       fake.CharactersN(10),
		Entry:             fake.CharactersN(20),
		Locale:            "en",
		InternalSignature: fake.CharactersN(15),
		CustomerID:        fake.CharactersN(8),
		DeliveryService:   fake.Company(),
		ShardKey:          fake.CharactersN(10),
		SMID:              fake.Day(),
		DateCreated:       "2021-11-26T06:22:19Z",
		OOFShard:          fake.CharactersN(5),
	}

	order.Delivery = types.Delivery{
		Name:    fake.FullName(),
		Phone:   fake.Phone(),
		Zip:     fake.Zip(),
		City:    fake.City(),
		Address: fake.StreetAddress(),
		Region:  fake.State(),
		Email:   fake.EmailAddress(),
	}

	order.Payment = types.Payment{
		Transaction:  fake.CharactersN(10),
		RequestID:    fake.CharactersN(8),
		Currency:     fake.CurrencyCode(),
		Provider:     fake.Company(),
		Amount:       3,
		PaymentDT:    time.Now().Unix(),
		Bank:         "Tinkoff",
		DeliveryCost: rand.Intn(10000),

		GoodsTotal: rand.Intn(100),
		CustomFee:  rand.Intn(10),
	}
	for i := 0; i < 3; i++ {
		item := types.Items{
			ChrtID:      i + 1,
			TrackNumber: fake.CharactersN(8),
			Price:       rand.Intn(30000),
			RID:         fake.CharactersN(5),
			Name:        fake.ProductName(),
			Sale:        rand.Intn(30),
			Size:        "0",
			TotalPrice:  rand.Intn(1000),
			NMID:        rand.Intn(599),
			Brand:       fake.Brand(),
			Status:      rand.Intn(400),
		}
		order.Items = append(order.Items, item)
	}
	// fmt.Println("orderUID of generated order ", order.OrderUID)
	return order
}
