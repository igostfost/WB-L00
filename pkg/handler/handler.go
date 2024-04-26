package handler

import (
	"WB_L00"
	"WB_L00/pkg/services"
	"WB_L00/types"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

type Handler struct {
	Service  *services.Service
	stanConn stan.Conn
	cache    *WB_L00.Cache
}

func NewHandler(service *services.Service, sc stan.Conn, cache *WB_L00.Cache) *Handler {
	return &Handler{
		Service:  service,
		stanConn: sc,
		cache:    cache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/*")
	router.GET("/", h.getForm)

	router.POST("/order", h.GetOrderByIDHandler)

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "notFound.html", nil)
	})
	return router
}

func (h *Handler) SubToChannel(channelName string) error {
	_, err := h.stanConn.Subscribe(channelName, func(msg *stan.Msg) {
		var newOrder types.Order
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Fatalf("error unmarshal order data: ", err)
			return
		}

		if !isValidNewOrder(newOrder) {
			log.Fatalf("invalid order data")
			return
		}

		err = h.Service.CreateOrder(newOrder)
		if err != nil {
			log.Fatalf("error create order: ", err)
		}

		h.cache.SetOrder(newOrder)
		log.Println("new order added in cache")

	}, stan.DurableName(channelName))
	if err != nil {
		return err
	}
	return nil
}

func isValidNewOrder(order types.Order) bool {
	if order.OrderUID == "" || order.TrackNumber == "" || order.CustomerID == "" {
		return false
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return false
	}
	if order.Payment.Transaction == "" {
		return false
	}
	for _, item := range order.Items {
		if item.ChrtID == 0 || item.TrackNumber == "" || item.Price == 0 || item.Name == "" {
			return false
		}
	}
	return true
}
