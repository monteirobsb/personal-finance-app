package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa o modelo de usuário no banco de dados
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Income        Income         `gorm:"foreignKey:UserID"`
	FixedExpenses []FixedExpense `gorm:"foreignKey:UserID"`
}

// AuthCode representa um código de autenticação enviado ao usuário
type AuthCode struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"index;not null"`
	CodeHash  string    `gorm:"not null"` // Código armazenado como hash
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

// Income representa a renda mensal do usuário
type Income struct {
	ID           uint    `gorm:"primaryKey"`
	UserID       uint    `gorm:"index;not null"` // Chave estrangeira para User
	MonthlyIncome float64 `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// FixedExpense representa uma despesa fixa mensal do usuário
type FixedExpense struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index;not null"` // Chave estrangeira para User
	Name      string  `gorm:"not null"`
	Value     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
