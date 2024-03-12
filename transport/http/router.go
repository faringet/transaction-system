package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"transaction-system/config"
)

type Router interface {
	Start()
	RegisterRoutes()
}

type RouterImpl struct {
	controller *Controller
	server     *gin.Engine
	logger     *zap.Logger
	url        string
}

func NewRouter(cfg *config.Config, logger *zap.Logger, watController *Controller) *RouterImpl {
	return &RouterImpl{controller: watController, logger: logger, url: cfg.LocalURL}
}

func (r *RouterImpl) RegisterRoutes() {
	router := gin.Default()

	router.POST("/invoice", func(c *gin.Context) {

		r.controller.AddAmount(c)
	})

	router.POST("/withdraw", func(c *gin.Context) {

		r.controller.WithdrawAmount(c)
	})

	router.GET("/available-balance", func(c *gin.Context) {

		r.controller.GetAvailableBalance(c)
	})

	router.GET("/frozen-balance", func(c *gin.Context) {

		r.controller.GetFrozenBalance(c)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "200 OK"})
	})

	r.server = router
}

func (r *RouterImpl) Start() error {
	return r.server.Run(r.url)
}
