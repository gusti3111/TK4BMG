package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah middleware yang memverifikasi token JWT dari header Authorization.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header diperlukan"})
			c.Abort()
			return
		}

		// 2. Cek format token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid (Gunakan Bearer)"})
			c.Abort()
			return
		}

		token := parts[1]

		// 3. Verifikasi Token (Skenario Produksi: Di sini seharusnya ada logika verifikasi JWT)
		// --- PLACEHOLDER LOGIC ---
		if token == "valid-bmg-token-123" {
			// Jika token valid, ekstrak user ID (misalnya dari klaim JWT)
			// Dalam contoh ini, kita hardcode user ID 1 sebagai user yang terautentikasi
			c.Set("user_id", 1) // Set User ID di Context untuk diakses oleh handler
			c.Next()            // Lanjutkan ke handler berikutnya
			return
		}

		// Jika token tidak valid
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
		c.Abort()
	}
}
