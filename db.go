package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Db struct {
	pool *pgxpool.Pool
}

func (db Db) List() []Product {
	rows, err := db.pool.Query(context.Background(), "SELECT * FROM products")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var products []Product

	for rows.Next() {
		p := Product{}
		if err := rows.Scan(&p.Name, &p.Category, &p.Price, &p.ID); err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products
}

func (db Db) Create(p Product) error {
	_, err := db.pool.Exec(context.Background(),
		"INSERT INTO products (name, category, price) VALUES($1, $2, $3)",
		p.Name, p.Category, p.Price)

	return err
}

func (db Db) Update(p Product) error {
	_, err := db.pool.Exec(context.Background(),
		"UPDATE products set name=$1, category=$2, price=$3 WHERE id=$4",
		p.Name, p.Category, p.Price, p.ID)

	return err
}
func (db Db) Delete(id int) error {
	_, err := db.pool.Exec(context.Background(), "DELETE from products where id = $1", id)
	return err
}
