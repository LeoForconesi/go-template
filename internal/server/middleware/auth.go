package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthConfig struct {
	Domain   string
	Audience string
}

func NewAuthMiddleware(cfg AuthConfig) (gin.HandlerFunc, error) {
	domain := strings.TrimSuffix(cfg.Domain, "/")
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	jwksURL := domain + "/.well-known/jwks.json"

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshUnknownKID: true,
	})
	if err != nil {
		return nil, err
	}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// 👇 OJO: en Auth0 "aud" puede ser string o array, esto es versión simple
		if aud, ok := claims["aud"].(string); !ok || aud != cfg.Audience {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid audience"})
			return
		}

		if iss, ok := claims["iss"].(string); !ok || iss != cfg.Domain {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
			return
		}

		// guardamos claims por si los querés usar después
		c.Set("user", claims)

		c.Next()
	}, nil
}
