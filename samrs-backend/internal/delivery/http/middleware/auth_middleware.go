package middleware

import (
	"fmt"
	"os"
	"strings"

	"samrs-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Pastikan sudah go get github.com/golang-jwt/jwt/v5
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.ErrorResponseFromErr(c, "Authorization header is required", util.ErrUnauthorized("authorization header missing"))
			c.Abort()
			return
		}

		// Format header: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.ErrorResponseFromErr(c, "Invalid authorization format", util.ErrUnauthorized("invalid authorization format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 2. Parse dan Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		// 3. Cek apakah token valid
		if err != nil || !token.Valid {
			util.ErrorResponseFromErr(c, "Invalid or expired token", util.ErrUnauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		// 4. Ekstrak data dari claims (payload token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			// Simpan data ke Context Gin agar bisa diakses di Handler/Controller
			c.Set("user_id", claims["user_id"])
			c.Set("tenant_id", claims["tenant_id"])
			c.Set("role_id", claims["role_id"])

			if userID, err := uuid.Parse(fmt.Sprint(claims["user_id"])); err == nil && userID != uuid.Nil {
				c.Set("user_id_uuid", userID)
			}
		}

		c.Next()
	}
}
