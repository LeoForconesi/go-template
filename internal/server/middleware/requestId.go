package middleware

/*
Qué hace: lee el header X-Request-ID; si no viene, genera un UUID y lo inyecta en request y response.
¿Para qué sirve? Trazabilidad: podés correlacionar logs de un request específico a través de todo el stack (API → otros servicios → colas). También ayuda cuando clientes reintentan llamadas: el mismo Request-ID viaja con el retry.
Extra: lo podés propagar a Kafka/Rabbit como header para atar eventos al request original.
*/
import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderRequestID = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderRequestID)
		if rid == "" {
			rid = uuid.NewString()
			c.Request.Header.Set(HeaderRequestID, rid)
		}
		c.Writer.Header().Set(HeaderRequestID, rid)
		c.Next()
	}
}
