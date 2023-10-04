package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
)

// 用于JWT签名和验证的密钥（通常应从安全存储中获取）
var jwtKey = []byte("your-secret-key")

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Protected resource")
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the old token from the request
	tokenAuth := jwtauth.New("HS256", []byte("your-secret-key"), nil)
	token, _, err := tokenAuth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get claims from the old token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Check if the token is about to expire
	if time.Now().Unix() > int64(claims["exp"].(float64))-60 { // Renew 60 seconds before expiration
		// Generate a new token with an extended expiration
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Extend token by 1 day

		newToken, _, err := tokenAuth.Encode(claims)
		if err != nil {
			http.Error(w, "Failed to renew token", http.StatusInternalServerError)
			return
		}

		// Send the new token to the client
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token": "` + newToken + `"}`))
		return
	}

	// Token is not expired; return the same token
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token": "` + token.Raw + `"}`))
}
