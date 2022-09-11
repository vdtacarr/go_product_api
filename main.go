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
		id, _ := ctx.ParamsInt("Id")
		err := db.Delete(id)
		if err != nil {
			ctx.SendString(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Listen(":8080")
}
