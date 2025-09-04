package account

import (
	"net/http"
	"go_starter_api/pkg/domain"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(accountService domain.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		accountID, err := accountService.ValidateAuthToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}

		c.Set("accountID", accountID)

		c.Next()
	}
}
