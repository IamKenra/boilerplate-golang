package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"svc-boilerplate-fiber-golang/utils/database"
	"svc-boilerplate-fiber-golang/utils/middleware"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	_masterdataHttpDeliver "svc-boilerplate-fiber-golang/domain/masterdata/delivery/http"
	_masterdataRepository "svc-boilerplate-fiber-golang/domain/masterdata/repository"
	_masterdataUsecase "svc-boilerplate-fiber-golang/domain/masterdata/usecase"
)

func main() {
	_ = godotenv.Load()
	if strings.TrimSpace(os.Getenv("JWT_SECRET")) == "" {
		log.Fatal("JWT_SECRET is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(middleware.CORSConfig())

	postgreDB := ConnectPostgre()
	if sqlDB, err := postgreDB.DB(); err == nil {
		defer sqlDB.Close()
	}

	masterdataRepo := _masterdataRepository.NewPostgresRepository(postgreDB)
	masterdataUC := _masterdataUsecase.NewMasterdataUsecase(masterdataRepo)
	_masterdataHttpDeliver.NewHttpMasterdataHandler(masterdataUC, app)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}

	listenErr := make(chan error, 1)
	go func() {
		listenErr <- app.Listen(":" + port)
	}()

	select {
	case <-ctx.Done():
		stop()
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	case err := <-listenErr:
		log.Fatal(err)
	}
}

func ConnectPostgre() *gorm.DB {
	postgre, err := database.SetupDatabaseConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	return postgre
}
