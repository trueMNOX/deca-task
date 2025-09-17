package user

import (
	"deca-task/internal/models"
	"deca-task/internal/user/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
func (h *userHandler) FindUsers(c *gin.Context) {
	var input struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
	input.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	input.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	phone := c.Query("phone")
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if phone != ""{
		query := h.userService.userRepository.db.Model(&models.User{})
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	users, total, err := h.userService.FindUsers(input.Page, input.Limit)
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
