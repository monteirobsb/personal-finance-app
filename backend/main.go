package main

import (
	"fmt"
	"log"
	// "net/http" // Gin vai cuidar disso
	"os"
	"personal-finance-app/backend/database"
	"personal-finance-app/backend/handlers"
	"personal-finance-app/backend/middleware" // Importa o pacote middleware

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http" // Necessário para http.StatusOK em helloHandler
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World from Backend with Gin!"})
}

func main() {
	// Carregar variáveis de ambiente de um arquivo .env (útil para desenvolvimento local)
	// Em produção, estas variáveis geralmente são configuradas diretamente no ambiente.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	// Definir as variáveis de ambiente para o banco de dados se não estiverem setadas
	// Isso é principalmente para o ambiente da sandbox, em um ambiente real
	// você configuraria isso fora da aplicação.
	setEnvIfNotExists("DB_HOST", "db") // Nome do serviço no docker-compose
	setEnvIfNotExists("DB_USER", "user")
	setEnvIfNotExists("DB_PASSWORD", "password")
	setEnvIfNotExists("DB_NAME", "personalfinancedb")
	setEnvIfNotExists("DB_PORT", "5432")


	// Conectar ao banco de dados
	database.ConnectDB()

	// Configurar o router Gin
	router := gin.Default()

	// Rota de exemplo
	router.GET("/", helloHandler)

	// Rotas de Autenticação
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/request-code", handlers.RequestCodeHandler)
		authRoutes.POST("/verify-code", handlers.VerifyCodeHandler)
	}

	// Rotas de Onboarding (protegidas por JWT)
	onboardingRoutes := router.Group("/onboarding")
	onboardingRoutes.Use(middleware.AuthMiddleware()) // Aplica o middleware de autenticação
	{
		onboardingRoutes.POST("/income", handlers.SaveIncomeHandler)
		onboardingRoutes.POST("/fixed-expenses", handlers.SaveFixedExpensesHandler)
	}

	// Rotas de Despesas Variáveis (protegidas por JWT)
	expenseRoutes := router.Group("/expenses")
	expenseRoutes.Use(middleware.AuthMiddleware())
	{
		expenseRoutes.POST("", handlers.PostExpenseHandler)       // POST /expenses
		expenseRoutes.DELETE("/:id", handlers.DeleteExpenseHandler) // DELETE /expenses/{id}
	}

	// Rota de Saldo e Projeção (protegida por JWT)
	router.GET("/balance", middleware.AuthMiddleware(), handlers.GetBalanceHandler)

	// Iniciar o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão se não especificada
	}
	fmt.Printf("Server started at port %s\n", port)
	log.Fatal(router.Run(":" + port))
}

// setEnvIfNotExists define uma variável de ambiente se ela ainda não estiver definida.
func setEnvIfNotExists(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}
