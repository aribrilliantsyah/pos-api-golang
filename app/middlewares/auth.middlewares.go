package middlewares

import (
	"context"
	"net/http"
	db "pos-api/db/sqlc"
	"pos-api/util/common"
	"pos-api/util/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(dbQ *db.Queries, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "authorization header is required",
			})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		userID, err := jwt.VerifyToken(bearerToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		userInfo, err := dbQ.GetUserByID(ctx, userID)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "token cannot be used (1)",
			})
			c.Abort()
			return
		}

		currentToken := common.ConvertNullString(userInfo.CurrentToken)
		if currentToken != bearerToken[1] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failed",
				"message": "token cannot be used (2)",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
