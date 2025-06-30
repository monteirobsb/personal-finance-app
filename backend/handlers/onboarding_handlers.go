package handlers

import (
	"log"
	"net/http"
	"personal-finance-app/backend/database"
	"personal-finance-app/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IncomePayload define a estrutura para receber a renda mensal
type IncomePayload struct {
	MonthlyIncome float64 `json:"rendaMensal" binding:"required,gte=0"`
}

// FixedExpensePayload define a estrutura para uma despesa fixa individual
type FixedExpensePayload struct {
	Name  string  `json:"nome" binding:"required"`
	Value float64 `json:"valor" binding:"required,gt=0"`
}

// FixedExpensesPayload define a estrutura para receber a lista de despesas fixas
type FixedExpensesPayload struct {
	Expenses []FixedExpensePayload `json:"despesasFixas" binding:"required,dive"`
}

// SaveIncomeHandler lida com o salvamento da renda mensal do usuário
func SaveIncomeHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var payload IncomePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Verificar se o usuário existe (embora o token JWT já deva garantir isso)
	var user models.User
	if result := database.DB.First(&user, uint(userID)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Salvar ou atualizar a renda
	// GORM `Assign` atualiza se existir (baseado em UserID), ou cria se não.
	// Precisamos garantir que o UserID está setado para o GORM encontrar ou criar corretamente.
	income := models.Income{
		UserID:        uint(userID),
		MonthlyIncome: payload.MonthlyIncome,
	}

	// Tentamos encontrar uma renda existente para este usuário
	var existingIncome models.Income
	result := database.DB.Where("user_id = ?", uint(userID)).First(&existingIncome)

	if result.Error == nil { // Renda existe, então atualizamos
		existingIncome.MonthlyIncome = payload.MonthlyIncome
		if err := database.DB.Save(&existingIncome).Error; err != nil {
			log.Printf("Error updating income for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save income"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Income updated successfully", "income": existingIncome})
	} else { // Renda não existe, criamos uma nova
		if err := database.DB.Create(&income).Error; err != nil {
			log.Printf("Error creating income for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save income"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Income saved successfully", "income": income})
	}
}

// SaveFixedExpensesHandler lida com o salvamento das despesas fixas do usuário
func SaveFixedExpensesHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var payload FixedExpensesPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Verificar se o usuário existe
	var user models.User
	if result := database.DB.First(&user, uint(userID)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Processar despesas: GORM não tem um "CreateOrUpdate" em lote fácil para sub-relações
	// A abordagem mais simples é deletar as antigas e criar as novas.
	// Outra abordagem seria iterar e atualizar/criar/deletar individualmente.
	// Para simplificar, vamos deletar as existentes e adicionar as novas.

	// Iniciar uma transação
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Error starting transaction for user %d: %v", userID, tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	// Deletar despesas fixas antigas para este usuário
	if err := tx.Where("user_id = ?", uint(userID)).Delete(&models.FixedExpense{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting old fixed expenses for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fixed expenses (delete step)"})
		return
	}

	// Criar novas despesas fixas
	var newExpenses []models.FixedExpense
	for _, expensePayload := range payload.Expenses {
		newExpenses = append(newExpenses, models.FixedExpense{
			UserID: uint(userID),
			Name:   expensePayload.Name,
			Value:  expensePayload.Value,
		})
	}

	if len(newExpenses) > 0 {
		if err := tx.Create(&newExpenses).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating new fixed expenses for user %d: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new fixed expenses"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error committing transaction for user %d: %v", userID, err.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit fixed expenses"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Fixed expenses saved successfully", "fixedExpenses": newExpenses})
}
