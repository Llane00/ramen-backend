package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config, env string) {
	var err error
	var dsn string

	switch env {
	case "development":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s_test port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	case "production":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	default:
		log.Fatal("Invalid environment specified")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the %s Database: %v", env, err)
	}
	fmt.Printf("✅ Connected Successfully to the %s Database\n", env)
}
