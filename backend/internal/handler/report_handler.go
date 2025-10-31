package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlerGenerateWeeklyReport(c *gin.Context) {
	// Implementasi di sini akan memanggil Report Service (area Rework)
	userID := c.GetInt("user_id") // Mengambil User ID dari AuthMiddleware
	c.JSON(http.StatusOK, gin.H{"message": "Generate Report logic placeholder. User ID:", "user_id": userID})
}
