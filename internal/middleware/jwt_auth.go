package middleware

import (
	"net/http"
	"strings"

	"github.com/prajwal-huggi/backend_go/internal/auth"
)

func JWTAuth(jwtService *auth.JWTService) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenStr := r.Header.Get("Authorization")

			if tokenStr == "" {
				http.Error(w, "missing token", 401)
				return
			}

			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			token, err := jwtService.ValidateToken(tokenStr)
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", 401)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}