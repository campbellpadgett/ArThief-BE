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
	page := c.Request.URL.Query().Get("page")
	// Get("page")
	if page == "" {
		c.Set("pageError", errors.New("pageID could not be retrieved"))
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.Set("pageError", err.Error())
	}

	if pageInt >= 21 {
		c.Set("pageError", errors.New("pageInt is too high"))
	}

	c.Set("pageInt", pageInt)

	userID := c.Request.URL.Query().Get("userID")
	if userID == "" {
		c.Set("userError", errors.New("userID could not be retrieved"))
	}

	c.Set("userID", userID)

	c.Next()
}
