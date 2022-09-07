package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	Name     string
	Category string
	Price    int
	ID       int
}

func main() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:12345@localhost:5432/postgres")
	if err != nil {
		fmt.Println("failed to connect to database: %w", err)
	}

	db := Db{pool: pool}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println(c.Method(), c.Path())
		return c.Next()
	})

	app.Get("/products", func(c *fiber.Ctx) error {
		return c.JSON(db.List())
	})

	app.Post("/products", func(c *fiber.Ctx) error {
		var product Product
		
		if err := c.BodyParser(&product); err != nil {
			return c.SendString(err.Error())
		}

		if err := db.Create(product); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	app.Put("/products/:id", func(ctx *fiber.Ctx) error {
		var p Product
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		id, _ := ctx.ParamsInt("id")
		p.ID = id

		err := db.Update(p)
		if err != nil {
			ctx.SendString(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})
    app.Delete("/deleteproduct/:id", func(ctx *fiber.Ctx) error {
		id := ctx.ParamsInt("Id")
		err := db.Delete(id)
		if err != nil {
			ctx.SendString(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Listen(":8080")
}

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
func (db Db) Delete(id int) error{
	_ , err := db.pool.Exec("DELETE from products where id = $1", id)
	return err
}
