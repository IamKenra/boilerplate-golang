package main

import (
	"database/sql"
	"log"
	"os"
	"boilerplate-golang/utils/database"
	"boilerplate-golang/utils/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	_masterdataHttpDeliver "boilerplate-golang/domain/masterdata/delivery/http"
	_masterdataRepository "boilerplate-golang/domain/masterdata/repository"
	_masterdataUsecase "boilerplate-golang/domain/masterdata/usecase"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()
	app.Use(middleware.CORSConfig())

	postgreDB := ConnectPostgre()

	masterdataRepo := _masterdataRepository.NewPostgresRepository(postgreDB)
	masterdataUC := _masterdataUsecase.NewMasterdataUsecase(masterdataRepo)
	_masterdataHttpDeliver.NewHttpMasterdataHandler(masterdataUC, app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
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
