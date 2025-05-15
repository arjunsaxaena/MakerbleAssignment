package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/arjunsaxaena/MakerbleAssignment/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret string

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found or error loading: %v", err)
	}
	
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("Warning: JWT_SECRET environment variable not set")
	} else {
		log.Println("JWT_SECRET successfully loaded")
	}
}

func GetJWTSecret() string {
	return jwtSecret
}

func GenerateToken(userID string, role string) (string, error) {
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userID,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   "Authorization header is required",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   "Authorization header must be in the format: Bearer {token}",
			})
			return
		}

		if jwtSecret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Error:   "JWT secret not configured",
			})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   "Invalid or expired token: " + err.Error(),
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["id"])
			c.Set("user_role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   "Invalid token claims",
			})
			return
		}
	}
}

func RoleAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Error:   "User role not found in token",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Error:   "Invalid role format in token",
			})
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, models.Response{
			Success: false,
			Error:   "You don't have permission to access this resource",
		})
	}
} 