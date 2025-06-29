package api

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		envPass := os.Getenv("TODO_PASSWORD")
		if envPass == "" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(envPass), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		if claims["pass_hash"] != hashPassword(envPass) {
			http.Error(w, "Invalid token content", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
