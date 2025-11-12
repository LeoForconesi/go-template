package middleware

/*
Timeout: envuelve el Context del request con un deadline (ej. 5s).
Si tu handler o sus llamadas bloquean, cortás la respuesta con 504. Evita colgar recursos. Útil con servicios externos.
*/
import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()
		select {
		case <-ctx.Done():
			c.AbortWithStatus(http.StatusGatewayTimeout)
		case <-done:
			return
		}
	}
}
