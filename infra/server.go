package infra

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port int
}

func NewServer(
	logger *logrus.Logger,
	config Config,
) *http.Server {
	router := gin.Default()

	rg := router.Group("/api/v1")

	rg.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return srv
}
