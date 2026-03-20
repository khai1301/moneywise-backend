package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParsePagination giúp parse tham số chia trang từ HTTP request, đồng thời giới hạn limit chống DDoS
func ParsePagination(c *gin.Context, maxLimit int) (limit int, offset int, page int) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	// Parse page
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Parse limit
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	// Hard cap for DoS prevention
	if limit > maxLimit {
		limit = maxLimit
	}

	offset = (page - 1) * limit
	return limit, offset, page
}
