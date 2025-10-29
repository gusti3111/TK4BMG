package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware mengatur header yang diperlukan agar permintaan cross-origin dari frontend diizinkan.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengizinkan semua origin (*) atau spesifikkan ke alamat frontend Anda (misal: "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// Metode yang diizinkan
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		// Header yang diizinkan untuk dikirim oleh client (sangat penting untuk Authorization)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Mengizinkan credentials (seperti cookies/session)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Jika ini adalah preflight request (OPTIONS), langsung kembalikan OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
