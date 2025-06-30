package handlers

import (
	"log"
	"math"
	"net/http"
	"personal-finance-app/backend/database"
	"personal-finance-app/backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BalanceResponse struct {
	CurrentBalance          float64     `json:"currentBalance"`
	TotalIncome             float64     `json:"totalIncome"`
	TotalFixedExpenses      float64     `json:"totalFixedExpenses"`
	TotalVariableExpenses   float64     `json:"totalVariableExpensesMonth"`
	Projection              *Projection `json:"projection,omitempty"`
	FinancialHealthStatus   string      `json:"financialHealthStatus"` // "verde", "amarelo", "vermelho"
	HealthPercentage        float64     `json:"healthPercentage"`
	DaysInMonthForProjection int        `json:"daysInMonthForProjection,omitempty"` // Para debug/info
	DayOfMonthForProjection  int        `json:"dayOfMonthForProjection,omitempty"`  // Para debug/info
}

type Projection struct {
	EndOfMonthBalance         float64 `json:"endOfMonthBalance"`
	ProjectedVariableExpenses float64 `json:"projectedVariableExpenses"`
	ProjectedTotalExpenses    float64 `json:"projectedTotalExpenses"`
	YellowAlertDay            string  `json:"yellowAlertDay,omitempty"` // Data "YYYY-MM-DD" ou dia do mês
	RedAlertDay               string  `json:"redAlertDay,omitempty"`    // Data "YYYY-MM-DD" ou dia do mês
	GMDVariableExpenses       float64 `json:"gmdVariableExpenses"`    // Gasto Médio Diário de Despesas Variáveis
}

// GetBalanceHandler calcula e retorna o saldo atual e a projeção.
func GetBalanceHandler(c *gin.Context) {
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

	// 1. Buscar Renda Mensal
	var income models.Income
	if err := database.DB.Where("user_id = ?", uint(userID)).First(&income).Error; err != nil {
		// Se não houver renda cadastrada, podemos retornar um erro ou um valor padrão.
		// Por enquanto, vamos assumir que o onboarding garantiu uma renda.
		// Em um app real, tratar o caso de não haver renda.
		log.Printf("Error fetching income for user %d: %v", userID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Income data not found for user. Please complete onboarding."})
		return
	}
	totalIncome := income.MonthlyIncome

	// 2. Buscar Total de Despesas Fixas
	var fixedExpenses []models.FixedExpense
	database.DB.Where("user_id = ?", uint(userID)).Find(&fixedExpenses)
	totalFixedExpenses := 0.0
	for _, fe := range fixedExpenses {
		totalFixedExpenses += fe.Value
	}

	// 3. Buscar Total de Despesas Variáveis no Mês Corrente
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1) // Último dia do mês corrente

	var variableExpensesMonth []models.VariableExpense
	database.DB.Where("user_id = ? AND date >= ? AND date <= ?", uint(userID), startOfMonth, now).Find(&variableExpensesMonth) // até o dia atual

	totalVariableExpensesMonth := 0.0
	for _, ve := range variableExpensesMonth {
		totalVariableExpensesMonth += ve.Value
	}

	// Calcular Saldo Atual
	currentBalance := totalIncome - totalFixedExpenses - totalVariableExpensesMonth

	// Calcular Saúde Financeira (para cor de fundo e status)
	healthPercentage := 0.0
	if totalIncome > 0 {
		// Net flow atual = Renda - Fixas - Variáveis do mês
		netFlowCurrent := totalIncome - totalFixedExpenses - totalVariableExpensesMonth
		healthPercentage = (netFlowCurrent / totalIncome) * 100
	}
	// Garante que não seja negativo para a lógica de cor do frontend que espera >= 0
	// healthPercentage = math.Max(healthPercentage, 0)

	financialHealthStatus := "vermelho" // Default
	if healthPercentage > 60 {
		financialHealthStatus = "verde"
	} else if healthPercentage >= 25 {
		financialHealthStatus = "amarelo"
	}


	var projectionData *Projection = nil
	daysInMonthForProjection := 0
	dayOfMonthForProjection := 0


	// Lógica de Projeção (se dia atual > 7 - ou seja, a partir do dia 8)
	dayOfMonth := now.Day()
	if dayOfMonth > 7 { // Condição: mais de 7 dias no mês (ou seja, a partir do dia 8)
		daysInMonth := float64(endOfMonth.Day()) // Número de dias no mês corrente
		daysInMonthForProjection = int(daysInMonth) // para debug
		dayOfMonthForProjection = dayOfMonth // para debug


		gmdVariableExpenses := 0.0
		if dayOfMonth > 0 && totalVariableExpensesMonth > 0 { // Evita divisão por zero se não houver gastos ou no primeiro dia
			gmdVariableExpenses = totalVariableExpensesMonth / float64(dayOfMonth)
		}

		projectedVariableExpensesMonth := gmdVariableExpenses * daysInMonth
		projectedTotalExpensesMonth := projectedVariableExpensesMonth + totalFixedExpenses
		projectedEndOfMonthBalance := totalIncome - projectedTotalExpensesMonth

		projectionData = &Projection{
			EndOfMonthBalance:         projectedEndOfMonthBalance,
			ProjectedVariableExpenses: projectedVariableExpensesMonth,
			ProjectedTotalExpenses:    projectedTotalExpensesMonth,
			GMDVariableExpenses:       gmdVariableExpenses,
		}

		// Estimativa de Dia para Alerta (Amarelo/Vermelho)
		if gmdVariableExpenses > 0 { // Só faz sentido projetar se houver um gasto médio diário
			currentSimBalance := totalIncome - totalFixedExpenses - totalVariableExpensesMonth

			foundYellow := false
			foundRed := false

			// Simula para os dias restantes no mês
			for d := dayOfMonth + 1; d <= int(daysInMonth); d++ {
				currentSimBalance -= gmdVariableExpenses

				// Calcula o percentual do saldo simulado em relação à renda
				// para determinar o estado (verde, amarelo, vermelho)
				// Saldo simulado aqui é o que sobraria da renda após todas as despesas fixas e variáveis até aquele dia.
				// Para fins de alerta, o que importa é o "dinheiro que sobra da renda mensal"
				// Então, o percentual é (currentSimBalance / totalIncome) * 100
				// No entanto, a lógica de cor é sobre o que *sobra* da renda após *todas* as despesas.
				// O saldo para alerta deve ser o saldo que o usuário teria no dia 'd'.
				// O currentSimBalance já é o saldo líquido projetado para o dia 'd'.
				// A questão é: o alerta é sobre o saldo absoluto ou sobre o percentual da renda que este saldo representa?
				// O enunciado diz: "dia estimado em que o saldo ficará 'amarelo' ou 'vermelho'".
				// As cores são definidas por faixas percentuais do (Receita - Despesas Fixas - Despesas Variáveis) / Receita.
				// Então, precisamos calcular esse percentual para o saldo simulado.

				// Saldo para o dia 'd' = Renda - Fixas - (Variáveis já ocorridas + Variáveis projetadas até 'd')
				// Variáveis projetadas até 'd' = gmd * (d - diaAtual)
				// Saldo no dia d = (Renda - Fixas - Variáveis_ja_ocorridas) - GMD * (d - diaAtual)
				//                = currentBalance_inicial_simulacao - GMD * (d - diaAtual)
				// currentSimBalance já é isso.

				simHealthPercentage := 0.0
				if totalIncome > 0 {
					simHealthPercentage = (currentSimBalance / totalIncome) * 100
				}
				// simHealthPercentage = math.Max(simHealthPercentage, 0) // Para consistência com a lógica de cor

				simDate := time.Date(now.Year(), now.Month(), d, 0, 0, 0, 0, now.Location())

				if !foundYellow && simHealthPercentage < 60 { // Limite para amarelo
					projectionData.YellowAlertDay = simDate.Format("2006-01-02")
					foundYellow = true
				}
				if !foundRed && simHealthPercentage < 25 { // Limite para vermelho
					projectionData.RedAlertDay = simDate.Format("2006-01-02")
					foundRed = true
				}
				if foundYellow && foundRed {
					break // Já encontrou ambos os alertas
				}
			}
		}
	}

	response := BalanceResponse{
		CurrentBalance:          currentBalance,
		TotalIncome:             totalIncome,
		TotalFixedExpenses:      totalFixedExpenses,
		TotalVariableExpenses:   totalVariableExpensesMonth,
		Projection:              projectionData,
		FinancialHealthStatus:   financialHealthStatus,
		HealthPercentage:        healthPercentage,
		DaysInMonthForProjection: daysInMonthForProjection,
		DayOfMonthForProjection:  dayOfMonthForProjection,
	}

	c.JSON(http.StatusOK, response)
}
