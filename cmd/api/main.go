package main

import (
	"fiberTest/config"
	"fiberTest/handlers"
	"fiberTest/log"
	"fiberTest/workerpool"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/uuid"
)

func main() {
	conf := config.New()
	var tasks []*workerpool.Task
	pool := workerpool.NewPool(tasks, conf.GetEnvInt("WORKERPOOL_SIZE", 5))

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Use(func(c fiber.Ctx) error {
		c.Locals("request_uuid", uuid.New().String())
		return c.Next()
	})

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/test", func(c fiber.Ctx) error {
		log.Info(c, "message", log.String("lol", "kek"), log.Int("val", 3))
		return c.SendString("Hello, World")
	})

	app.Get("/request/:id", func(c fiber.Ctx) error {
		return handlers.GetRequest(c, conf)
	})

	app.Post("/request", func(c fiber.Ctx) error {
		return handlers.CreateRequest(c, pool, conf)
	})

	go func() {
		pool.RunBackground()
	}()

	app.Listen(conf.GetEnv("APP_PORT", ":3000") /*, fiber.ListenConfig{EnablePrefork: true}*/)
}
