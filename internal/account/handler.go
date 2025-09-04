package account

import (
	"net/http"
	"servicehub_api/pkg/domain"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	logger            *logrus.Logger
	accountService    domain.AccountService
	accountRepository domain.AccountRepository
}

func NewAccountHandler(
	logger *logrus.Logger,
	accountService domain.AccountService,
	accountRepository domain.AccountRepository,
) *AccountHandler {
	return &AccountHandler{
		logger:            logger,
		accountService:    accountService,
		accountRepository: accountRepository,
	}
}

type RegisterAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterAccountResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// @Summary		Register a new account
// @Description	Register a new account
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			account	body		RegisterAccountRequest	true	"Account"
// @Success		200		{object}	RegisterAccountResponse
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			/account/register [post]
func (h *AccountHandler) RegisterAccount(c *gin.Context) {
	var req RegisterAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if account already exists
	existingAcc, err := h.accountRepository.GetAccountByEmail(req.Email)
	if err == nil && existingAcc != nil {
		h.logger.Errorf("account already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "account already exists"})
		return
	}
	if err != nil {
		h.logger.Errorf("failed to get account by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Hash the password before storing
	hashedPassword, err := h.accountService.HashPassword(req.Password)
	if err != nil {
		h.logger.Errorf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	acc := &domain.Account{
		Email:    req.Email,
		Password: hashedPassword,
	}

	acc, err = h.accountRepository.CreateAccount(acc)
	if err != nil {
		h.logger.Errorf("failed to create account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.accountService.GenerateToken(acc)
	if err != nil {
		h.logger.Errorf("failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RegisterAccountResponse{
		ID:    acc.ID,
		Email: acc.Email,
		Token: token,
	})
}
