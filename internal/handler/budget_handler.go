package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/config"
	"github.com/khai1301/moneywise-backend/internal/service"
)

type BudgetHandler struct {
	budgetService service.BudgetService
}

func NewBudgetHandler(budgetService service.BudgetService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

type BudgetRequest struct {
	CategoryID string  `json:"categoryId" binding:"required,uuid4"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Month      string  `json:"month" binding:"required"` // "YYYY-MM"
	Note       string  `json:"note" binding:"omitempty,max=500"`
}

func (h *BudgetHandler) Create(c *gin.Context) {
	var req BudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}
	userID := c.MustGet("user_id").(string)

	budget, err := h.budgetService.CreateBudget(userID, req.CategoryID, req.Month, req.Note, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tạo ngân sách thành công", "budget": budget})
}

func (h *BudgetHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	month := c.DefaultQuery("month", time.Now().Format("2006-01"))

	budgets, err := h.budgetService.GetBudgetsWithStats(userID, month, config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy ngân sách"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": budgets, "month": month})
}

func (h *BudgetHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	var req BudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	budget, err := h.budgetService.UpdateBudget(id, userID, req.CategoryID, req.Month, req.Note, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công", "budget": budget})
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	if err := h.budgetService.DeleteBudget(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa ngân sách thành công"})
}
