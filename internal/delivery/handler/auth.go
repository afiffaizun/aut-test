package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"auth-service/internal/domain"
	"auth-service/internal/pkg/response"
	"auth-service/internal/usecase"
)

type AuthHandler struct {
	registerUseCase  *usecase.RegisterUseCase
	loginUseCase    *usecase.LoginUseCase
	refreshUseCase  *usecase.RefreshUseCase
	logoutUseCase   *usecase.LogoutUseCase
	validateUseCase *usecase.ValidateUseCase
}

func NewAuthHandler(
	register *usecase.RegisterUseCase,
	login *usecase.LoginUseCase,
	refresh *usecase.RefreshUseCase,
	logout *usecase.LogoutUseCase,
	validate *usecase.ValidateUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase:  register,
		loginUseCase:    login,
		refreshUseCase:  refresh,
		logoutUseCase:   logout,
		validateUseCase: validate,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinBadRequest(c, err)
		return
	}

	if err := validator.New().Struct(req); err != nil {
		response.GinBadRequest(c, err)
		return
	}

	user, tokens, err := h.registerUseCase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == usecase.ErrEmailAlreadyExists {
			response.GinBadRequest(c, err)
			return
		}
		response.GinInternalServerError(c, err)
		return
	}

	response.GinCreated(c, gin.H{
		"user":   domain.UserResponse{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt},
		"tokens": tokens,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinBadRequest(c, err)
		return
	}

	user, tokens, err := h.loginUseCase.Execute(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			response.GinUnauthorized(c, err)
			return
		}
		response.GinInternalServerError(c, err)
		return
	}

	response.GinOK(c, gin.H{
		"user":   domain.UserResponse{ID: user.ID, Email: user.Email, CreatedAt: user.CreatedAt},
		"tokens": tokens,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.GinBadRequest(c, err)
		return
	}

	tokens, err := h.refreshUseCase.Execute(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == domain.ErrInvalidToken || err == domain.ErrTokenExpired || err == domain.ErrTokenRevoked {
			response.GinUnauthorized(c, err)
			return
		}
		response.GinInternalServerError(c, err)
		return
	}

	response.GinOK(c, gin.H{"tokens": tokens})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	accessTokenID := c.GetHeader("X-Token-ID")

	expiry := 15 * time.Minute

	err := h.logoutUseCase.Execute(c.Request.Context(), refreshToken, accessTokenID, expiry)
	if err != nil {
		response.GinInternalServerError(c, err)
		return
	}

	response.GinOK(c, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		response.GinUnauthorized(c, domain.ErrInvalidToken)
		return
	}

	response.GinOK(c, gin.H{"user_id": userID})
}

func (h *AuthHandler) Validate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != "" {
		token = token[7:] // Remove "Bearer " prefix
	}

	claims, err := h.validateUseCase.Execute(c.Request.Context(), token)
	if err != nil {
		response.GinUnauthorized(c, err)
		return
	}

	response.GinOK(c, gin.H{"claims": claims})
}