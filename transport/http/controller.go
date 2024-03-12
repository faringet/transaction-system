package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"transaction-system/internal/domain"
)

type Wat interface {
	AddAmountController(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error)
	WithdrawAmountController(c *gin.Context, currencyCode int, amount float64, walletNumber int, cardNumber int) (*domain.Transactions, error)
	GetAvailableBalanceController(c *gin.Context, walletNumber int, cardNumber int) (float64, error)
	GetFrozenBalanceController(c *gin.Context, walletNumber int, cardNumber int) (float64, error)
}

type Controller struct {
	wat    Wat
	logger *zap.Logger
}

func NewWatController(wat Wat, logger *zap.Logger) *Controller {
	return &Controller{wat: wat, logger: logger}
}

func (c2 *Controller) AddAmount(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c2.logger.Error("Failed to parse request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	transaction, err := c2.wat.AddAmountController(c, req.CurrencyCode, req.Amount, req.WalletNumber, req.CardNumber)
	if err != nil {
		c2.logger.Error("Failed to add amount to database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add amount to database"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (c2 *Controller) WithdrawAmount(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c2.logger.Error("Failed to parse request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	transaction, err := c2.wat.WithdrawAmountController(c, req.CurrencyCode, req.Amount, req.WalletNumber, req.CardNumber)
	if err != nil {
		c2.logger.Error("Failed to withdraw amount from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to withdraw amount from database"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (c2 *Controller) GetAvailableBalance(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c2.logger.Error("Failed to parse request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	availableBalance, err := c2.wat.GetAvailableBalanceController(c, req.WalletNumber, req.CardNumber)
	if err != nil {
		c2.logger.Error("Failed to fetch available balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_balance": availableBalance})
}

func (c2 *Controller) GetFrozenBalance(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c2.logger.Error("Failed to parse request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	frozenBalance, err := c2.wat.GetFrozenBalanceController(c, req.WalletNumber, req.CardNumber)
	if err != nil {
		c2.logger.Error("Failed to fetch frozen balance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch frozen balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"frozen_balance": frozenBalance})
}
