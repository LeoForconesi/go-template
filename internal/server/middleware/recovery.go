package middleware

/*
captura panics en handlers para que el proceso no caiga; devuelve 500 y loguea el stacktrace.Es un “Airbag” del servidor. Onda si explota algo que no esta contemplado
*/
import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered", zap.Any("error", rec))
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
