package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"svc-boilerplate-fiber-golang/utils/database"
	"svc-boilerplate-fiber-golang/utils/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	_masterdataHttpDeliver "svc-boilerplate-fiber-golang/domain/masterdata/delivery/http"
	_masterdataRepository "svc-boilerplate-fiber-golang/domain/masterdata/repository"
	_masterdataUsecase "svc-boilerplate-fiber-golang/domain/masterdata/usecase"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()
	app.Use(middleware.CORSConfig())

	postgreDB := ConnectPostgre()

	masterdataRepo := _masterdataRepository.NewPostgresRepository(postgreDB)
	masterdataUC := _masterdataUsecase.NewMasterdataUsecase(masterdataRepo)
	_masterdataHttpDeliver.NewHttpMasterdataHandler(masterdataUC, app)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func ConnectPostgre() (postgre *sql.DB) {
	postgre, err := database.SetupPostgreDatabaseConnection()

	if err != nil {
		log.Fatal(err.Error())
	}

	return
}

type metricsJobUsecase interface {
	RunHourlyMetricsJob() error
}
