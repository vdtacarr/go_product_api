package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"go_db/entities"
	"go_db/services"
)

func main() {
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/pegasus_app")
	if err != nil {
		fmt.Println("failed to connect to database: %w", err)
	}

	db := services.Db{Pool: pool}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println(c.Method(), c.Path())
		return c.Next()
	})
	app.Get("/products", func(c *fiber.Ctx) error {
		return c.JSON(db.List())
	})
	app.Get("products/:id", func(c *fiber.Ctx) error {
		id, _ := c.ParamsInt("id")
		product := db.GetById(id)
		if product.Id == 0 {
			c.SendString("No Item Found!")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(product)
	})
	app.Post("/products", func(c *fiber.Ctx) error {
		var product entities.Product

		if err := c.BodyParser(&product); err != nil {
			return c.SendString(err.Error())
		}

		if err := db.Create(product); err != nil {
			c.SendString(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusCreated)
	})
	app.Put("/update-product/:id", func(ctx *fiber.Ctx) error {
		var p entities.Product
		if err := ctx.BodyParser(&p); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		id, _ := ctx.ParamsInt("id")
		p.Id = id

		err := db.Update(p)
		if err != nil {
			ctx.SendString(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Delete("/delete-product/:id", func(ctx *fiber.Ctx) error {
		id, _ := ctx.ParamsInt("Id")
		err := db.Delete(id)
		if err != nil {
			ctx.SendString(err.Error())
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})
	app.Listen(":3000")
}
