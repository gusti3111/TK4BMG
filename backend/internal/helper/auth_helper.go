package helper

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserID mengambil userID dari context yang diisi oleh middleware JWT.
// Fungsi ini bisa dipakai ulang di semua handler.
func GetUserID(c *gin.Context) (int, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		log.Println("[Helper] UserID tidak ditemukan di context (JWT mungkin tidak valid)")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak terotentikasi"})
		return 0, false
	}

	userID, ok := userIDValue.(int)
	if !ok {
		log.Println("[Helper] Format UserID tidak valid di context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format UserID tidak valid"})
		return 0, false
	}

	return userID, true
}
