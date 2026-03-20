package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/service"
	"github.com/khai1301/moneywise-backend/pkg/utils"
)

type TransactionHandler struct {
	txService service.TransactionService
}

func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

type TransactionRequest struct {
	CategoryID    string  `json:"categoryId" binding:"required,uuid4"`
	Title         string  `json:"title" binding:"required,min=1,max=255"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Type          string  `json:"type" binding:"required,oneof=income expense transfer"`
	Date          string  `json:"date" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" binding:"omitempty,max=50"`
	Note          string  `json:"note" binding:"omitempty,max=1000"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ (Dùng chuẩn RFC3339 ví dụ 2024-01-01T00:00:00Z)"})
		return
	}

	tx, err := h.txService.CreateTransaction(
		userID, req.CategoryID, req.Title, req.Amount, req.Type, date, req.PaymentMethod, req.Note,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Thêm giao dịch thành công", "transaction": tx})
}

func (h *TransactionHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	// Sử dụng Hàm tiện ích chống DDoS với hardcap = 100 limit max
	limit, offset, page := utils.ParsePagination(c, 100)

	var startDate, endDate time.Time
	if start := c.Query("start_date"); start != "" {
		startDate, _ = time.Parse(time.RFC3339, start)
	}
	if end := c.Query("end_date"); end != "" {
		endDate, _ = time.Parse(time.RFC3339, end)
	}

	txs, total, err := h.txService.GetTransactionsByUser(userID, startDate, endDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lấy danh sách giao dịch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": txs,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	tx, err := h.txService.GetTransactionByID(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if tx == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy giao dịch"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": tx})
}

func (h *TransactionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ: " + err.Error()})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng ngày không hợp lệ"})
		return
	}

	tx, err := h.txService.UpdateTransaction(
		id, userID, req.CategoryID, req.Title, req.Amount, req.Type, date, req.PaymentMethod, req.Note,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật thành công", "transaction": tx})
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("user_id").(string)

	err := h.txService.DeleteTransaction(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Xóa thành công"})
}
