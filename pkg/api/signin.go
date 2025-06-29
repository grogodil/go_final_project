package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	envPass := os.Getenv("TODO_PASSWORD")
	if envPass == "" {
		writeJSONError(w, "аутентификация не настроена", http.StatusInternalServerError)
		return
	}

	if req.Password != envPass {
		writeJSONError(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	// Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pass_hash": hashPassword(envPass),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(envPass))
	if err != nil {
		writeJSONError(w, "ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]string{"token": tokenString}, http.StatusOK)
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
