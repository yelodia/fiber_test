package handlers

import (
	"fiberTest/config"
	"fiberTest/log"
	"fiberTest/models"
	"fiberTest/workerpool"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"time"
)

func GetRequest(c fiber.Ctx, conf *config.Config) error {
	id := c.Params("id")
	log.Info(c, "get request", log.String("id", id))
	row := conf.DB.QueryRow("SELECT * FROM requests WHERE uuid = $1", id)

	var req models.Request
	err := row.Scan(&req.UUID, &req.Value, &req.Timeout, &req.Result, &req.State)
	if err != nil {
		log.Error(c, "Parsing error", log.Err(err))
		return c.Status(404).SendString("Request not found")
	}

	return c.JSON(req)
}

func CreateRequest(c fiber.Ctx, pool *workerpool.Pool, conf *config.Config) error {
	req := new(models.Request)
	log.Info(c, "start request", log.String("data", string(c.Body())))
	if err := c.Bind().Body(req); err != nil {
		log.Error(c, "Parsing error", log.Err(err))
		return c.Status(400).SendString("Invalid request")
	}

	req.UUID = uuid.New().String()

	_, err := conf.DB.Exec("INSERT INTO requests (uuid, value, timeout) VALUES ($1, $2, $3)",
		req.UUID, req.Value, req.Timeout)
	if err != nil {
		log.Error(c, "Request create error", log.Err(err))
		return c.Status(500).SendString("Request create error")
	}

	task := workerpool.NewTask(func(data interface{}) error {
		r := data.(*models.Request)
		_, e := conf.DB.Exec("UPDATE requests SET state = $1 WHERE uuid = $2",
			models.RequestStatePending, r.UUID)
		if e != nil {
			log.Error(c, "Request update error", log.Err(e))
			return e
		}
		time.Sleep(time.Duration(r.Timeout) * time.Millisecond)
		r.Result = r.Value * r.Value
		if r.Result == 0 {
			r.State = models.RequestStateError
		} else {
			r.State = models.RequeststateSuccess
		}
		_, e = conf.DB.Exec("UPDATE requests SET state = $1, result = $2 WHERE uuid = $3",
			r.State, r.Result, r.UUID)
		if e != nil {
			log.Error(c, "Request update error", log.Err(e))
			return e
		}
		return nil
	}, req)
	pool.AddTask(task)

	return c.Status(201).SendString(req.UUID)

}
