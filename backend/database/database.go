package database

import (
	"fmt"
	"log"
	"os"
	"personal-finance-app/backend/models" // Ajuste o path se necessário

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB inicializa a conexão com o banco de dados PostgreSQL
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully opened")

	// Migrar o schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.AuthCode{},
		&models.Income{},
		&models.FixedExpense{},
		&models.VariableExpense{}, // Adiciona VariableExpense à migração
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migrated")
}
