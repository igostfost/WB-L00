package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getForm(c *gin.Context) {
	c.HTML(http.StatusOK, "searchForm.html", nil)
}

func (h *Handler) GetOrderByIDHandler(c *gin.Context) {
	orderID := c.PostForm("id")
	order, found := h.cache.GetOrderByUID(orderID)
	if !found {
		c.Redirect(http.StatusFound, "/?not_found")
		return
	}

	c.HTML(http.StatusOK, "orderForm.html", gin.H{
		"OrderUID":          order.OrderUID,
		"TrackNumber":       order.TrackNumber,
		"Entry":             order.Entry,
		"Delivery":          order.Delivery,
		"Payment":           order.Payment,
		"OrderItems":        order.Items,
		"Locale":            order.Locale,
		"InternalSignature": order.InternalSignature,
		"CustomerID":        order.CustomerID,
		"DeliveryService":   order.DeliveryService,
		"ShardKey":          order.ShardKey,
		"SMID":              order.SMID,
		"DateCreated":       order.DateCreated,
		"OOFShard":          order.OOFShard,
	})
}
