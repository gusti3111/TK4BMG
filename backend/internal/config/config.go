package config

import "os"

// JWTSecretKey adalah kunci rahasia global untuk JWT.
// Diambil dari environment variable untuk keamanan,
// dengan fallback ke nilai default jika tidak diset.
var JWTSecretKey = getJWTSecret()

func getJWTSecret() []byte {
	// Best practice: Ambil secret dari environment variable
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		// Fallback untuk development jika env var tidak diset
		// JANGAN GUNAKAN INI DI PRODUKSI
		return []byte("your-very-secret-key")
	}
	return []byte(secret)
}
