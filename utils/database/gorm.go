package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabaseConnection() (*gorm.DB, error) {
	dialect := strings.ToLower(strings.TrimSpace(os.Getenv("DB_DIALECT")))
	if dialect == "" {
		dialect = "postgres"
	}

	var (
		dialector gorm.Dialector
		name      string
	)
	switch dialect {
	case "postgres", "postgresql":
		dialector = postgres.Open(postgresDSN())
		name = "PostgreSQL"
	case "mysql":
		dialector = mysql.Open(mysqlDSN())
		name = "MySQL"
	default:
		return nil, fmt.Errorf("unsupported DB_DIALECT: %s", dialect)
	}

	if missing := findMissingEnvValues(map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USERNAME": os.Getenv("DB_USERNAME"),
	}); len(missing) > 0 {
		return nil, fmt.Errorf("missing required database environment variables: %s", strings.Join(missing, ", "))
	}

	gormLogger := logger.Default.LogMode(logger.Silent)
	if os.Getenv("DB_DEBUG") == "true" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open %s connection: %w", name, err)
	}
	if db == nil {
		return nil, errors.New("failed to open database connection: nil database handle")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s sql handle: %w", name, err)
	}

	sqlDB.SetMaxOpenConns(envInt("DB_MAX_OPEN_CONNS", 25))
	sqlDB.SetMaxIdleConns(envInt("DB_MAX_IDLE_CONNS", 25))
	sqlDB.SetConnMaxLifetime(envDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute))
	sqlDB.SetConnMaxIdleTime(envDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("failed to ping %s connection: %w", name, err)
	}

	return db, nil
}

func postgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func mysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FJakarta",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}
