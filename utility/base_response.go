package utility

import "github.com/gin-gonic/gin"

type (
	// ResponseType ...
	ResponseType int
	// ErrResponse ...
	ErrResponse struct {
		Message string
	}
)

const (
	// JSON ...
	JSON ResponseType = iota
	// XML ...
	XML
)

// RouteResponseType ...
var RouteResponseType ResponseType

// ResponseWithType ...
func ResponseWithType(c *gin.Context, statusCode int, body interface{}) {
	switch RouteResponseType {
	case XML:
		c.Abort()
		c.XML(statusCode, body)
	case JSON:
		fallthrough
	default:
		c.AbortWithStatusJSON(statusCode, body)
	}
}
