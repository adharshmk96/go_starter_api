package infra

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Config struct {
	Port int
}

func NewServer(
	db *gorm.DB,
	logger *logrus.Logger,
	config Config,
) *http.Server {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	rg := router.Group("/api/v1")

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	SetupRoutes(rg, db, logger)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return srv
}
