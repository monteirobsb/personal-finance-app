module personal-finance-app/backend

go 1.21 // Ou a versão Go que você pretende usar

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/joho/godotenv v1.5.1 // Adicionado para carregar .env
	golang.org/x/crypto v0.17.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)
