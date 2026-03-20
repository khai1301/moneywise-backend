package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	user, err := h.userService.UpdateProfile(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request or passwords too short"})
		return
	}

	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
