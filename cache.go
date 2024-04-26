package WB_L00

import (
	"WB_L00/types"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]types.Order
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]types.Order)}
}

func (c *Cache) SetOrder(order types.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUID] = order
}

func (c *Cache) GetOrderByUID(uid string) (types.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, ok := c.data[uid]
	return order, ok
}

func (c *Cache) RestoreFromDB(orders []types.Order) {
	for _, order := range orders {
		c.data[order.OrderUID] = order
	}
}

func (c *Cache) CleanCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]types.Order)
}
