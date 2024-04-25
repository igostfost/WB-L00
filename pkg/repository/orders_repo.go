package repository

import (
	"WB_L00/types"
	"fmt"
)

//func CreateOrder

//func GetAllOrders

//func GetOrderByUID

func (r *Repository) CreateOrder(newOrder types.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	orderQwery := `INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(orderQwery, newOrder.OrderUID,
		newOrder.TrackNumber,
		newOrder.Entry,
		newOrder.Locale,
		newOrder.InternalSignature,
		newOrder.CustomerID,
		newOrder.DeliveryService,
		newOrder.ShardKey,
		newOrder.SMID,
		newOrder.DateCreated,
		newOrder.OOFShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	deliveryQwery := `INSERT INTO delivery (order_uid, name, phone, zip, city,address,region, email) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = tx.Exec(deliveryQwery, newOrder.OrderUID,
		newOrder.Delivery.Name,
		newOrder.Delivery.Phone,
		newOrder.Delivery.Zip,
		newOrder.Delivery.City,
		newOrder.Delivery.Address,
		newOrder.Delivery.Region,
		newOrder.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	paymentQwery := `INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err = tx.Exec(paymentQwery, newOrder.OrderUID,
		newOrder.Payment.Transaction,
		newOrder.Payment.RequestID,
		newOrder.Payment.Currency,
		newOrder.Payment.Provider,
		newOrder.Payment.Amount,
		newOrder.Payment.PaymentDT,
		newOrder.Payment.Bank,
		newOrder.Payment.DeliveryCost,
		newOrder.Payment.GoodsTotal,
		newOrder.Payment.CustomFee)

	if err != nil {
		tx.Rollback()
		return err
	}

	itemsQwery := `INSERT INTO order_items (order_uid, track_number, chrt_id, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range newOrder.Items {
		_, err = tx.Exec(itemsQwery, newOrder.OrderUID,
			item.TrackNumber,
			item.ChrtID,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			item.Brand,
			item.Status)

		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) GetAllOrdersFromDB() ([]types.Order, error) {
	query := `
        SELECT
            o.order_uid,
            o.track_number,
            o.entry,
            o.locale,
            o.internal_signature,
            o.customer_id,
            o.delivery_service,
            o.shardkey,
            o.sm_id,
            o.date_created,
            o.oof_shard,
            d.name AS delivery_name,
            d.phone AS delivery_phone,
            d.zip AS delivery_zip,
            d.city AS delivery_city,
            d.address AS delivery_address,
            d.region AS delivery_region,
            d.email AS delivery_email,
            p.transaction AS payment_transaction,
            p.request_id AS payment_request_id,
            p.currency AS payment_currency,
            p.provider AS payment_provider,
            p.amount AS payment_amount,
            p.payment_dt AS payment_payment_dt,
            p.bank AS payment_bank,
            p.delivery_cost AS payment_delivery_cost,
            p.goods_total AS payment_goods_total,
            p.custom_fee AS payment_custom_fee
        FROM
            orders o
        LEFT JOIN
            delivery d ON o.order_uid = d.order_uid
        LEFT JOIN
            payment p ON o.order_uid = p.order_uid
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var orders []types.Order

	for rows.Next() {
		var order types.Order
		var delivery types.Delivery
		var payment types.Payment

		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.DateCreated,
			&order.OOFShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.RequestID,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		order.Delivery = delivery
		order.Payment = payment

		orderItems, err := r.getOrderItemsByOrderUID(order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order items: %v", err)
		}
		order.Items = orderItems

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %v", err)
	}

	return orders, nil
}

func (r *Repository) getOrderItemsByOrderUID(orderUID string) ([]types.Items, error) {
	query := `
		SELECT
			oi.chrt_id AS item_chrt_id,
			oi.track_number AS item_track_number,
			oi.price AS item_price,
			oi.rid AS item_rid,
			oi.name AS item_name,
			oi.sale AS item_sale,
			oi.size AS item_size,
			oi.total_price AS item_total_price,
			oi.nm_id AS item_nm_id,
			oi.brand AS item_brand,
			oi.status AS item_status
		FROM
			order_items oi
		WHERE
			oi.order_uid = $1
	`

	// Выполнить запрос SQL
	rows, err := r.db.Query(query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var orderItems []types.Items

	for rows.Next() {
		var orderItem types.Items

		err := rows.Scan(
			&orderItem.ChrtID,
			&orderItem.TrackNumber,
			&orderItem.Price,
			&orderItem.RID,
			&orderItem.Name,
			&orderItem.Sale,
			&orderItem.Size,
			&orderItem.TotalPrice,
			&orderItem.NMID,
			&orderItem.Brand,
			&orderItem.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %v", err)
	}

	return orderItems, nil
}
