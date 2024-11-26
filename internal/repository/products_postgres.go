package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	products "github.com/lavatee/shop_products"
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
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

func (r *ProductsPostgres) DeleteProduct(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", productsTable)
	_, err := r.db.Exec(query, id)
	return err
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
	query := fmt.Sprintf("SELECT FROM %s WHERE id=$1", productsTable)
	err := r.db.Get(&product, query, id)
	return product, err
}

func (r *ProductsPostgres) PostOrder(productId int) error {
	query := fmt.Sprintf("UPDATE %s SET amount=amount-1 WHERE id=$1", productsTable)
	_, err := r.db.Exec(query, productId)
	return err
}
