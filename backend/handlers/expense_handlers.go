package handlers

import (
	"net/http"
	"personal-finance-app/backend/database"
	"personal-finance-app/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateExpensePayload define a estrutura para criar uma nova despesa variável
type CreateExpensePayload struct {
	Value       float64 `json:"value" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Description string  `json:"description"` // Opcional
	Date        string  `json:"date"`        // Opcional, formato "YYYY-MM-DD"
}

// PostExpenseHandler lida com o registro de uma nova despesa variável
func PostExpenseHandler(c *gin.Context) {
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

	var payload CreateExpensePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	expenseDate := time.Now() // Default para hoje
	if payload.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", payload.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
			return
		}
		expenseDate = parsedDate
	}

	variableExpense := models.VariableExpense{
		UserID:      uint(userID),
		Value:       payload.Value,
		Category:    payload.Category,
		Description: payload.Description,
		Date:        expenseDate,
	}

	if result := database.DB.Create(&variableExpense); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save expense: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Expense registered successfully", "expense": variableExpense})
}

// DeleteExpenseHandler lida com a remoção de uma despesa variável
func DeleteExpenseHandler(c *gin.Context) {
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

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID format"})
		return
	}

	var expense models.VariableExpense
	// Primeiro, encontrar a despesa para garantir que pertence ao usuário e existe
	if result := database.DB.Where("id = ? AND user_id = ?", uint(expenseID), uint(userID)).First(&expense); result.Error != nil {
		if result.Error.Error() == "record not found" { // gorm.ErrRecordNotFound não é exportado diretamente em alguns contextos
			c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found or you do not have permission to delete it."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + result.Error.Error()})
		}
		return
	}

	// Se encontrada e pertence ao usuário, deletar
	if result := database.DB.Delete(&expense); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense: " + result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
