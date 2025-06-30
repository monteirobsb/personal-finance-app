package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"personal-finance-app/backend/database"
	"personal-finance-app/backend/models"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Deve ser configurado via variável de ambiente

// EmailRegex é uma regex simples para validação de e-mail.
// Para uma validação robusta em produção, considere bibliotecas especializadas.
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// RequestCodeBody define a estrutura esperada para o corpo da requisição /auth/request-code
type RequestCodeBody struct {
	Email string `json:"email" binding:"required"`
}

// VerifyCodeBody define a estrutura esperada para o corpo da requisição /auth/verify-code
type VerifyCodeBody struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// generateSecureCode gera um código alfanumérico seguro de n dígitos.
func generateSecureCode(length int) (string, error) {
	bytes := make([]byte, length/2+1) // Cada byte em hexadecimal são 2 caracteres
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := hex.EncodeToString(bytes)
	return strings.ToUpper(code[:length]), nil // Garante o comprimento e caixa alta
}

// hashData gera um hash bcrypt de uma string.
func hashData(data string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// checkDataHash compara uma string com seu hash bcrypt.
func checkDataHash(data, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

// RequestCodeHandler lida com a solicitação de um código de autenticação.
func RequestCodeHandler(c *gin.Context) {
	var body RequestCodeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validar formato do e-mail
	if !EmailRegex.MatchString(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Gerar código de 6 dígitos
	code, err := generateSecureCode(6)
	if err != nil {
		log.Printf("Error generating secure code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication code"})
		return
	}

	// Hashear o código
	codeHash, err := hashData(code)
	if err != nil {
		log.Printf("Error hashing code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure authentication code"})
		return
	}

	// Salvar no banco de dados
	authCodeEntry := models.AuthCode{
		Email:     body.Email,
		CodeHash:  codeHash,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if result := database.DB.Create(&authCodeEntry); result.Error != nil {
		log.Printf("Error saving auth code to DB: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save authentication code"})
		return
	}

	// Simular envio de e-mail (log por enquanto)
	log.Printf("Auth code for %s: %s (Simulated email send)", body.Email, code)

	c.JSON(http.StatusOK, gin.H{"message": "Authentication code sent (simulated)."})
}

// VerifyCodeHandler lida com a verificação do código e login/criação de usuário.
func VerifyCodeHandler(c *gin.Context) {
	var body VerifyCodeBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Validar formato do e-mail
	if !EmailRegex.MatchString(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}
	if len(body.Code) != 6 { // Assumindo que o código tem sempre 6 dígitos
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code format"})
		return
	}

	// Buscar código no banco de dados
	var authCodeEntry models.AuthCode
	if result := database.DB.Where("email = ?", body.Email).Order("created_at desc").First(&authCodeEntry); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or code. Code not found."})
		return
	}

	// Verificar se o código expirou
	if time.Now().After(authCodeEntry.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication code expired."})
		return
	}

	// Verificar o hash do código
	if !checkDataHash(body.Code, authCodeEntry.CodeHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication code."})
		return
	}

	// Código válido, encontrar ou criar usuário
	var user models.User
	if result := database.DB.Where("email = ?", body.Email).FirstOrCreate(&user, models.User{Email: body.Email}); result.Error != nil {
		log.Printf("Error finding or creating user: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user account"})
		return
	}

	// Gerar token JWT
	// Definir a chave secreta JWT (deve ser carregada de forma segura, ex: variável de ambiente)
	if jwtKey == nil || len(jwtKey) == 0 {
		// Tenta carregar do ambiente se ainda não estiver setada (ex: no início da aplicação)
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Println("CRITICAL: JWT_SECRET is not set. Using a default insecure key for now. THIS IS NOT SAFE FOR PRODUCTION.")
			// Para ambientes de teste/desenvolvimento sem a variável setada, podemos usar um default.
			// **NUNCA FAÇA ISSO EM PRODUÇÃO.**
			jwtKey = []byte("default_insecure_secret_key_for_testing_only_12345")
		} else {
			jwtKey = []byte(secret)
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24 horas
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprint(user.ID), // Armazena o ID do usuário no token
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		return
	}

	// Opcionalmente, invalidar o código de autenticação após o uso bem-sucedido
	// database.DB.Delete(&authCodeEntry) // Ou marcar como usado

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully authenticated.",
		"token":   tokenString,
		"userId":  user.ID,
		"email":   user.Email,
	})
}
