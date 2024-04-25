package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func (h *Handler) GetOrderByIDHandler(c *gin.Context) {
	orderID := c.PostForm("id")
	fmt.Println(orderID)
	order, found := h.cache.GetOrderByUID(orderID)
	if !found {
		c.Redirect(http.StatusFound, "/?not_found")
		return
	}

	// Передаем данные заказа в HTML-шаблон и отображаем его
	//c.HTML(http.StatusOK, "order.html", gin.H{
	//	"OrderUID":          order.OrderUID,
	//	"TrackNumber":       order.TrackNumber,
	//	"Entry":             order.Entry,
	//	"Delivery":          order.Delivery,
	//	"Payment":           order.Payment,
	//	"OrderItems":        order.Items,
	//	"Locale":            order.Locale,
	//	"InternalSignature": order.InternalSignature,
	//	"CustomerID":        order.CustomerID,
	//	"DeliveryService":   order.DeliveryService,
	//	"ShardKey":          order.ShardKey,
	//	"SMID":              order.SMID,
	//	"DateCreated":       order.DateCreated,
	//	"OOFShard":          order.OOFShard,
	//})

	// Выводим данные заказа в консоль
	fmt.Println("Order Details:")
	fmt.Println("OrderUID:", order.OrderUID)
	fmt.Println("TrackNumber:", order.TrackNumber)
	fmt.Println("Entry:", order.Entry)
	fmt.Println("Delivery:", order.Delivery)
	fmt.Println("Payment:", order.Payment)
	fmt.Println("Locale:", order.Locale)
	fmt.Println("InternalSignature:", order.InternalSignature)
	fmt.Println("CustomerID:", order.CustomerID)
	fmt.Println("DeliveryService:", order.DeliveryService)
	fmt.Println("ShardKey:", order.ShardKey)
	fmt.Println("SMID:", order.SMID)
	fmt.Println("DateCreated:", order.DateCreated)
	fmt.Println("OOFShard:", order.OOFShard)

	// Возвращаем сообщение об успешном завершении
	fmt.Println("Order details printed successfully!")
}
