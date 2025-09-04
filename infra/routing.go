package infra

import (
	"servicehub_api/internal/account"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SetupRoutes(rg *gin.RouterGroup, db *gorm.DB, logger *logrus.Logger) {
	accountRepository := account.NewAccountRepository(db)
	accountService := account.NewAccountService()
	accountHandler := account.NewAccountHandler(logger, accountService, accountRepository)

	rg.POST("/account/register", accountHandler.RegisterAccount)
}
