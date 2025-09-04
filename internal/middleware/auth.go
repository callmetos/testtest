package middleware

import (
	"net/http"
	"strings"

	"navmate-backend/pkg/jwtauth"

	"github.com/gin-gonic/gin"
)

func AuthJWT(jwt *jwtauth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(h[len("Bearer "):])
		claims, err := jwt.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", int(claims.UserID))
		c.Set("email", claims.Email)
		c.Next()
	}
}
