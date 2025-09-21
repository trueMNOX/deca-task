package user

import (
	"deca-task/internal/user/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type UserListResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int64          `json:"total"`
	Users []UserResponse `json:"users"`
}

type userHandler struct {
	userService *userService
}

func NewUserHandler(userService *userService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) UsersRoute(r *gin.RouterGroup) {
	r.GET("/users/:id", h.FindUserById)
	r.GET("/users", h.FindUsers)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v2/users/{id} [get]
func (h *userHandler) FindUserById(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		return
	}

	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	user, err := h.userService.FindUserById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID,
		Phone:     user.PhoneNumber,
		CreatedAt: user.CreatedAt,
	})
}

// GetUsers godoc
// @Summary List users
// @Description Get users with pagination and optional phone search
// @Tags user
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param phone query string false "Search by phone"
// @Success 200 {object} UserListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v2/users [get]
func (h *userHandler) FindUsers(c *gin.Context) {
	var input struct {
        Page  int    `form:"page"`
        Limit int    `form:"limit"`
        Phone string `form:"phone"`
    }
	input.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	input.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	input.Phone = c.Query("phone")

	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, total, err := h.userService.FindUsers(input.Page, input.Limit, input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Phone:     user.PhoneNumber,
			CreatedAt: user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.UserListResponse{
		Page:  input.Page,
		Limit: input.Limit,
		Total: int64(total),
		Users: userResponses,
	})
}
