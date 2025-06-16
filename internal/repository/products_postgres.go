package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	products "github.com/lavatee/shop_products"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

const (
	deletedStatus   = "deleted"
	availableStatus = "available"
)

type ProductsPostgres struct {
	db    *sqlx.DB
	cache *redis.Client
}

func NewProductsPostgres(db *sqlx.DB, cache *redis.Client) *ProductsPostgres {
	return &ProductsPostgres{
		db: db,
	}
}

func (r *ProductsPostgres) PostProduct(name string, amount int, price int, category string, description string, userId int) (int, error) {
	var id int
	fmt.Println(amount)
	query := fmt.Sprintf("INSERT INTO %s (name, amount, price, category, description, user_id) values ($1, $2, $3, $4, $5, $6) RETURNING id", productsTable)
	row := r.db.QueryRow(query, name, amount, price, category, description, userId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductsPostgres) PostDeleteProductEvent(id int, eventId string, productCreator int) error {
	if !r.CheckIsUserProduct(id, productCreator) {
		return fmt.Errorf("That's not a user's product")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}
	if err := tx.Rollback(); err != nil {
		return err
	}
	query = fmt.Sprintf("UPDATE %s SET status = $1, status_event_id = $2", productsTable)
	_, err = tx.Exec(query, deletedStatus, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductsPostgres) CheckIsUserProduct(productId int, userId int) bool {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE id = $1 AND user_id = $2")
	row := r.db.QueryRow(query, productId, userId)
	if err := row.Scan(&id); err != nil {
		return false
	}
	return true
}

func (r *ProductsPostgres) GetProducts(category string) ([]products.Product, error) {
	query := fmt.Sprintf("SELECT FROM %s WHERE category=$1", productsTable)
	var products []products.Product
	if err := r.db.Select(&products, query, category); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsPostgres) GetUserProducts(userId int) ([]products.Product, error) {
	query := fmt.Sprintf("SELECT FROM %s WHERE user_id=$1", productsTable)
	var products []products.Product
	if err := r.db.Select(&products, query, userId); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsPostgres) GetSavedProducts(ids []int) ([]products.Product, error) {
	query := fmt.Sprintf("SELECT FROM %s WHERE", productsTable)
	for index := range ids {
		query += fmt.Sprintf(" id=$%x", index+1)
		if index+1 < len(ids) {
			query += " OR"
		}
	}
	args := make([]interface{}, len(ids))
	for index, id := range ids {
		args[index] = id
	}
	var products []products.Product
	if err := r.db.Select(&products, query, args...); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsPostgres) GetOneProduct(id int) (products.Product, error) {
	var product products.Product
	cachedProduct, err := r.cache.Get(ctx, fmt.Sprint(id)).Result()
	if err != redis.Nil {
		query := fmt.Sprintf("SELECT FROM %s WHERE id=$1", productsTable)
		err := r.db.Get(&product, query, id)
		if err != nil {
			return products.Product{}, err
		}
		jsonProduct, err := json.Marshal(product)
		if err != nil {
			return products.Product{}, err
		}
		if err := r.cache.Set(ctx, fmt.Sprint("id"), jsonProduct, cachedItemTTL).Err(); err != nil {
			return products.Product{}, err
		}
		return product, nil
	}
	if err := json.Unmarshal([]byte(cachedProduct), &product); err != nil {
		return products.Product{}, err
	}
	return product, nil
}

func (r *ProductsPostgres) PostEvent(eventId string, eventType string, confirmersAmount int) error {
	query := fmt.Sprintf("INSERT INTO %s (mq_event_id, type, confirmers_amount, confirmations_amount) VALUES ($1, $2, $3)", eventsTable)
	_, err := r.db.Exec(query, eventId, eventType, confirmersAmount, 0)
	return err
}

func (r *ProductsPostgres) PostCompensatoryEvent(comEventId string, eventId string) error {
	var compensatedEventId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE mq_event_id = $1", eventsTable)
	row := r.db.QueryRow(query, eventId)
	if err := row.Scan(&compensatedEventId); err != nil {
		return err
	}
	query = fmt.Sprintf("INSERT INTO %s (mq_event_id, event_id) VALUES ($1, $2)", compensatoryEventsTable)
	_, err := r.db.Exec(query, comEventId, compensatedEventId)
	return err
}

func (r *ProductsPostgres) CompensateDeleteProductEvent(eventId string) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE status_event_id = $2", productsTable)
	_, err := r.db.Exec(query, availableStatus, eventId)
	return err
}

type ConfirmEventInfo struct {
	ConfirmationsAmount int
	ConfirmersAmount    int
}

func (r *ProductsPostgres) ConfirmProductDeleting(eventId string) error {
	var info ConfirmEventInfo
	query := fmt.Sprintf("UPDATE %s SET confirmations_amount = confirmations_amount + 1 WHERE mq_event_id = $1 RETURNING confirmations_amount, confirmers_amount", eventsTable)
	if err := r.db.Get(&info, query, eventId); err != nil {
		return err
	}
	if info.ConfirmationsAmount != info.ConfirmersAmount {
		return nil
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE status = $1 AND status_event_id = $2", productsTable)
	_, err := r.db.Exec(query, deletedStatus, eventId)
	return err
}

func (r *ProductsPostgres)