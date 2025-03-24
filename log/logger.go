package log

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	fmt.Println("init logger")
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err) // Не удалось создать логгер
	}
	defer logger.Sync()
}

func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Err(err error) zap.Field {
	return zap.Error(err)
}

func Info(fctx fiber.Ctx, msg string, fields ...zap.Field) {
	fields = append(fields, zap.Any("uuid", fctx.Locals("request_uuid")))
	logger.Info(msg,
		fields...,
	)
}

func Error(fctx fiber.Ctx, msg string, fields ...zap.Field) {
	fields = append(fields, zap.Any("uuid", fctx.Locals("request_uuid")))
	logger.Error(msg,
		fields...,
	)
}
