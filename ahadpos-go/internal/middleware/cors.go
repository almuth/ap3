package middleware

import (
    "strings"

    "github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(allowedOrigins []string, allowedMethods []string, allowedHeaders []string, maxAge int) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // Check if origin is allowed
        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if allowedOrigin == "*" || allowedOrigin == origin {
                allowed = true
                break
            }
        }

        if allowed {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
        }

        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
        c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
        c.Writer.Header().Set("Access-Control-Max-Age", string(rune(maxAge)))

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
