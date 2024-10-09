package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Llane00/ramen-backend/initializers"
	"github.com/Llane00/ramen-backend/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("❌ Could not load environment variables", err)
	}

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // default env
	}

	initializers.ConnectDB(&config, env)
}

func main() {
	fmt.Println("GO_ENV:", os.Getenv("GO_ENV"))
	initializers.DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Shop{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)
	fmt.Println("✅ Migration complete")
}
