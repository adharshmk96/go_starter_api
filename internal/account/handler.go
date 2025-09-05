package account

import (
	"errors"
	"net/http"
	"servicehub_api/pkg/domain"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
// @Router			api/v1/account/register [post]
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Errorf("failed to get account by email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
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

	token, err := h.accountService.GenerateAuthToken(acc)
	if err != nil {
		h.logger.Errorf("failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.accountRepository.LogAccountActivity(acc.ID, domain.ActivityRegister)
	if err != nil {
		h.logger.Errorf("failed to log activity: %v", err)
	}

	c.JSON(http.StatusOK, RegisterAccountResponse{
		ID:    acc.ID,
		Email: acc.Email,
		Token: token,
	})
}

type LoginAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAccountResponse struct {
	Token string `json:"token"`
}

// @Summary		Login a user
// @Description	Login a user
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			account	body		LoginAccountRequest	true	"Account"
// @Success		200		{object}	LoginAccountResponse
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			api/v1/account/login [post]
func (h *AccountHandler) LoginAccount(c *gin.Context) {
	var req LoginAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.accountRepository.GetAccountByEmail(req.Email)
	if err != nil {
		h.logger.Errorf("failed to get account by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ok, err := h.accountService.ComparePassword(req.Password, acc.Password)
	if err != nil {
		h.logger.Errorf("failed to compare password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if !ok {
		h.logger.Errorf("invalid password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	token, err := h.accountService.GenerateAuthToken(acc)
	if err != nil {
		h.logger.Errorf("failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	err = h.accountRepository.LogAccountActivity(acc.ID, domain.ActivityLogin)
	if err != nil {
		h.logger.Errorf("failed to log activity: %v", err)
	}

	c.JSON(http.StatusOK, LoginAccountResponse{Token: token})
}

// @Summary		Logout a user
// @Description	Logout a user
// @Tags			account
// @Accept			json
// @Produce		json
// @Success		200		{object}	map[string]string
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			api/v1/account/logout [post]
func (h *AccountHandler) LogoutAccount(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	accountID, err := h.accountService.ValidateAuthToken(token)
	if err != nil {
		h.logger.Errorf("failed to validate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err = h.accountRepository.LogAccountActivity(accountID, domain.ActivityLogout)
	if err != nil {
		h.logger.Errorf("failed to log activity: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

type GetProfileResponse struct {
	ID         uint                     `json:"id"`
	Email      string                   `json:"email"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	Activities []domain.AccountActivity `json:"activities"`
}

// @Summary		Get Profile
// @Description	Get Profile of the authenticated user
// @Tags			account
// @Accept			json
// @Produce		json
// @Success		200		{object}	ProfileResponse
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			api/v1/account/profile [get]
func (h *AccountHandler) GetProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	accountID, err := h.accountService.ValidateAuthToken(token)
	if err != nil {
		h.logger.Errorf("failed to validate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	acc, err := h.accountRepository.GetAccountByID(accountID)
	if err != nil {
		h.logger.Errorf("failed to get account by id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, GetProfileResponse{ID: acc.ID, Email: acc.Email})
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// @Summary		Forgot Password
// @Description	Forgot Password
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			account	body		ForgotPasswordRequest	true	"Account"
// @Success		200		{object}	ForgotPasswordResponse
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			api/v1/account/forgot-password [post]
func (h *AccountHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acc, err := h.accountRepository.GetAccountByEmail(req.Email)
	if err != nil {
		h.logger.Errorf("failed to get account by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if acc == nil {
		h.logger.Errorf("account not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		return
	}

	token, err := h.accountService.GeneratePasswordResetToken(acc)
	if err != nil {
		h.logger.Errorf("failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	err = h.accountService.SendPasswordResetEmail(acc.Email, token)
	if err != nil {
		h.logger.Errorf("failed to send password reset email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send password reset email"})
		return
	}

	err = h.accountRepository.LogAccountActivity(acc.ID, domain.ActivityForgotPassword)
	if err != nil {
		h.logger.Errorf("failed to log activity: %v", err)
	}

	c.JSON(http.StatusOK, ForgotPasswordResponse{Message: "password reset email sent"})
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// @Summary		Reset Password
// @Description	Reset Password
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			account	body		ResetPasswordRequest	true	"Account"
// @Success		200		{object}	ResetPasswordResponse
// @Failure		400		{object}	map[string]string
// @Failure		500		{object}	map[string]string
// @Router			api/v1/account/reset-password [post]
func (h *AccountHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := req.Token
	password := req.Password

	accountID, err := h.accountService.ValidatePasswordResetToken(token)
	if err != nil {
		h.logger.Errorf("failed to validate token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	acc, err := h.accountRepository.GetAccountByID(accountID)
	if err != nil {
		h.logger.Errorf("failed to get account by id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	hashedPassword, err := h.accountService.HashPassword(password)
	if err != nil {
		h.logger.Errorf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	acc.Password = hashedPassword

	acc, err = h.accountRepository.UpdateAccount(acc)
	if err != nil {
		h.logger.Errorf("failed to update account: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err = h.accountRepository.LogAccountActivity(acc.ID, domain.ActivityResetPassword)
	if err != nil {
		h.logger.Errorf("failed to log activity: %v", err)
	}

	c.JSON(http.StatusOK, ResetPasswordResponse{Message: "password reset successful"})
}
