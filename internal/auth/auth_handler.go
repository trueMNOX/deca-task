package auth

import (
	"deca-task/internal/auth/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService *authservice
}

func NewAuthHandler(authservice *authservice) *authHandler {
	return &authHandler{authService: authservice}
}
func (h *authHandler) AuthRoute(r *gin.RouterGroup) {
	r.POST("/login", h.Login)
	r.POST("/verify", h.VerifyOTP)
}
func (h *authHandler) Login(c *gin.Context) {
	var input dto.RequestOTPDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	otp, err := h.authService.LoginUser(input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.RequestOTPResponse{
		Message: "OTP sent successfully",
		OTP:     otp.OTP,
	})
}
func (h *authHandler) VerifyOTP(c *gin.Context) {
	var input dto.VerifyOTPDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := h.authService.VerifyOTP(input.Phone, input.OTP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &dto.VerifyOTPResponse{Token: user.Token})
}
