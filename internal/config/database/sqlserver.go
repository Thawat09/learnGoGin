package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectSQLServer(host, port, user, password string) (*gorm.DB, error) {

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=testITD", user, password, host, port)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 5 * time.Second,
				LogLevel:      logger.Warn,
				Colorful:      false,
			},
		),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	DB = db

	return db, nil
}
