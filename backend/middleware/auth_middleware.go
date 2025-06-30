package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware é um middleware para verificar o token JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inicializa jwtKey se estiver vazia (fallback para desenvolvimento, como em auth_handlers)
		if jwtKey == nil || len(jwtKey) == 0 {
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				log.Println("CRITICAL (AuthMiddleware): JWT_SECRET is not set. Using a default insecure key. THIS IS NOT SAFE FOR PRODUCTION.")
				jwtKey = []byte("default_insecure_secret_key_for_testing_only_12345")
			} else {
				jwtKey = []byte(secret)
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// O token geralmente vem no formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]

		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verifica o método de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			validationErr, ok := err.(*jwt.ValidationError)
			if ok {
				if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Malformed token"})
					return
				} else if validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is expired or not valid yet"})
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Token é válido. Armazenar o ID do usuário (Subject do token) no contexto do Gin.
		// O ID do usuário foi armazenado no campo Subject do RegisteredClaims.
		c.Set("userID", claims.Subject)

		c.Next()
	}
}
