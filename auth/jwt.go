package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// 用于JWT签名和验证的密钥（通常应从安全存储中获取）
var jwtKey = []byte("your-secret-key")

// Token验证中间件
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取Bearer令牌
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// 解析JWT令牌
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// 如果令牌验证通过，继续处理请求
		next.ServeHTTP(w, r)
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Protected resource")
}
