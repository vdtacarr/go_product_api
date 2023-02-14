package services

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go_db/entities"
)

type Db struct {
	Pool *pgxpool.Pool
}

func (db Db) List() []entities.Product {
	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var products []entities.Product

	for rows.Next() {
		p := entities.Product{}
		if err := rows.Scan(&p.Id, &p.Name, &p.Category, &p.Price); err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products
}
func (db Db) Create(p entities.Product) error {
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO products (name, category, price) VALUES($1, $2, $3)",
		p.Name, p.Category, p.Price)

	return err
}
func (db Db) Update(p entities.Product) error {
	_, err := db.Pool.Exec(context.Background(),
		"UPDATE products set name=$1, category=$2, price=$3 WHERE id=$4",
		p.Name, p.Category, p.Price, p.Id)

	return err
}
func (db Db) Delete(id int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE from products where id = $1", id)
	return err
}
func (db Db) GetById(id int) entities.Product {
	product := entities.Product{}
	row := db.Pool.QueryRow(context.Background(), "SELECT * FROM products where id = $1", id)
	row.Scan(&product.Id, &product.Name, &product.Category, &product.Price)
	return product
}
