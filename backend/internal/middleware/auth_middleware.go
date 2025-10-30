package middleware

import (
	"fmt" // <-- DITAMBAHKAN
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Kunci rahasia untuk JWT.
// PENTING: Kunci ini HARUS SAMA PERSIS dengan kunci yang Anda gunakan di `auth_handler.go`!
// (Sebaiknya simpan ini di file .env dan impor dari config)
var jwtSecretKey = []byte("kunci_rahasia_anda_yang_sangat_aman")

// AuthMiddleware adalah middleware yang memvalidasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Auth Error: Header Authorization kosong")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Permintaan tidak terotentikasi (header kosong)"})
			return
		}

		// 2. Periksa format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Auth Error: Format header tidak valid (bukan Bearer)")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Format token tidak valid"})
			return
		}

		tokenString := parts[1]

		// 3. Parse dan Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing (HS256) sudah benar
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// --- PERBAIKAN DI SINI ---
				errMsg := fmt.Sprintf("Metode signing tidak terduga: %v", token.Header["alg"])
				log.Println(errMsg)
				return nil, fmt.Errorf(errMsg) // Mengembalikan error yang valid
			}
			// Kembalikan secret key (HARUS SAMA DENGAN SAAT SIGNING)
			return jwtSecretKey, nil
		})

		// 4. Tangani Error Validasi (Ini adalah tempat error Anda terjadi)
		if err != nil {
			log.Printf("Auth Error: Token tidak valid: %v", err)
			// Ini adalah respons yang dilihat frontend Anda
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kedaluwarsa"})
			return
		}

		// 5. Ekstrak Claims (Payload)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			// Ambil 'sub' (User ID) dan 'role' dari token
			// (Sesuai dengan apa yang kita masukkan di `auth_handler.go`)
			userIDFloat, okSub := claims["sub"].(float64) // JWT menyimpan angka sebagai float64
			userRole, okRole := claims["role"].(string)

			if !okSub || !okRole {
				log.Println("Auth Error: Claims 'sub' (userID) atau 'role' tidak ditemukan/invalid di token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Payload token tidak valid"})
				return
			}

			// Konversi userID (sub) ke int
			userID := int(userIDFloat)

			// 6. SIMPAN DATA USER DI CONTEXT GIN
			// Ini agar handler selanjutnya (seperti ItemHandler) bisa mengakses siapa yang login
			c.Set("userID", userID)
			c.Set("userRole", userRole)

			// Lanjutkan ke handler berikutnya
			c.Next()

		} else {
			log.Println("Auth Error: Claims token tidak valid atau token tidak aktif")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims token tidak valid"})
			return
		}
	}
}
