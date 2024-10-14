// middleware/authorize.go

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizeRoles(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized"})
			return
		}

		role, ok := userRole.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "User role type assertion failed"})
			return
		}

		for _, allowedRole := range roles {
			if role == allowedRole {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Access denied"})
	}
}
