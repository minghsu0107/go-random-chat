package common

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTTPContextKey string

var (
	JWTAuthHeader                = "Authorization"
	ChannelKey    HTTPContextKey = "channel_key"
)

func MaxAllowed(n int64) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()

	}
}

// CORSMiddleware adds CORS headers to each response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := extractTokenFromHeader(c.Request)
		if accessToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authResult, err := Auth(&AuthPayload{
			AccessToken: accessToken,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if authResult.Expired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrResponse{
				Message: ErrTokenExpired.Error(),
			})
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ChannelKey, authResult.ChannelID))
		c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearToken := r.Header.Get(JWTAuthHeader)
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
