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

// RequestOTP godoc
// @Summary Request OTP
// @Description Send phone number to receive an OTP (valid for 2 minutes)
// @Tags auth
// @Accept  json
// @Produce  json
// @Param data body dto.RequestOTPDTO true "Request body"
// @Success 200 {object} dto.RequestOTPResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/login [post]
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

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP and return JWT if valid
// @Tags auth
// @Accept  json
// @Produce  json
// @Param data body dto.VerifyOTPDTO true "Request body"
// @Success 200 {object} dto.VerifyOTPResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/verify [post]
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
