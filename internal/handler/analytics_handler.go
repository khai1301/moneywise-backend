package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khai1301/moneywise-backend/internal/config"
	"github.com/khai1301/moneywise-backend/internal/models"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

// GET /api/analytics/monthly?year=2026
// Returns 12 months of income, expense, savings for the given year.
func (h *AnalyticsHandler) Monthly(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	yearStr := c.DefaultQuery("year", fmt.Sprintf("%d", time.Now().Year()))
	year := yearStr

	type MonthStat struct {
		Month   string  `json:"month"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
		Savings float64 `json:"savings"`
	}

	months := []MonthStat{}
	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	for m := 1; m <= 12; m++ {
		monthStr := fmt.Sprintf("%s-%02d", year, m)
		start, _ := time.Parse("2006-01", monthStr)
		end := start.AddDate(0, 1, 0).Add(-time.Second)

		var income, expense float64
		config.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND type = 'income' AND date >= ? AND date <= ?", userID, start, end).
			Select("COALESCE(SUM(amount), 0)").Scan(&income)
		config.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND type = 'expense' AND date >= ? AND date <= ?", userID, start, end).
			Select("COALESCE(SUM(amount), 0)").Scan(&expense)

		months = append(months, MonthStat{
			Month:   monthNames[m-1],
			Income:  income,
			Expense: expense,
			Savings: income - expense,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": months, "year": year})
}

// GET /api/analytics/categories?start=YYYY-MM-DD&end=YYYY-MM-DD
// Returns spending totals grouped by category.
func (h *AnalyticsHandler) CategorySummary(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	now := time.Now()
	// Default: current month
	defaultStart := fmt.Sprintf("%d-%02d-01", now.Year(), int(now.Month()))
	defaultEnd := fmt.Sprintf("%d-%02d-%02d", now.Year(), int(now.Month()), daysInMonth(now))

	startStr := c.DefaultQuery("start", defaultStart)
	endStr   := c.DefaultQuery("end", defaultEnd)

	start, _ := time.Parse("2006-01-02", startStr)
	end, _   := time.Parse("2006-01-02", endStr)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	type CatRow struct {
		CategoryID   string  `json:"categoryId"`
		CategoryName string  `json:"name"`
		Icon         string  `json:"icon"`
		Color        string  `json:"color"`
		Total        float64 `json:"value"`
	}

	var rows []CatRow
	config.DB.Table("transactions t").
		Select("t.category_id, c.name as category_name, c.icon, c.color, COALESCE(SUM(t.amount), 0) as total").
		Joins("JOIN categories c ON c.id = t.category_id").
		Where("t.user_id = ? AND t.type = 'expense' AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL", userID, start, end).
		Group("t.category_id, c.name, c.icon, c.color").
		Order("total DESC").
		Scan(&rows)

	// Compute total for percentages
	var grandTotal float64
	for _, r := range rows {
		grandTotal += r.Total
	}

	type CatResult struct {
		CategoryID string  `json:"categoryId"`
		Name       string  `json:"name"`
		Icon       string  `json:"icon"`
		Color      string  `json:"color"`
		Value      float64 `json:"value"`
		Pct        float64 `json:"pct"`
	}
	results := make([]CatResult, 0, len(rows))
	for _, r := range rows {
		pct := 0.0
		if grandTotal > 0 {
			pct = (r.Total / grandTotal) * 100
		}
		results = append(results, CatResult{
			CategoryID: r.CategoryID,
			Name:       r.CategoryName,
			Icon:       r.Icon,
			Color:      r.Color,
			Value:      r.Total,
			Pct:        pct,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        results,
		"grandTotal":  grandTotal,
		"start":       startStr,
		"end":         endStr,
	})
}

func daysInMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
