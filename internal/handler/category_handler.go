package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/service"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

type CategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Type        string `json:"type" binding:"required,oneof=income expense both"`
	Icon        string `json:"icon" binding:"max=255"`
	Color       string `json:"color" binding:"max=50"`
	Description string `json:"description" binding:"max=1000"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)

	category, err := h.categoryService.CreateCategory(userID, req.Name, req.Type, req.Icon, req.Color, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tạo danh mục thành công", "category": category})
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	categories, err := h.categoryService.GetCategoriesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách danh mục"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	category, err := h.categoryService.UpdateCategory(id, userID, req.Name, req.Type, req.Icon, req.Color, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công", "category": category})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	err := h.categoryService.DeleteCategory(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
