package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func CorsMiddleware(origins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Paginate(c *gin.Context) {
	page, ok := c.Params.Get("page")
	if !ok {
		c.Set("error", errors.New("pageID could not be retrieved"))
	}

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			c.Set("error", err.Error())
		}

		c.Set("pageInt", pageInt)
	}

	c.Next()
}
