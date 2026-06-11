package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

			claims := token.Claims.(jwt.MapClaims)
			fmt.Println("JWT Claims:", claims)
			ctx := context.WithValue(r.Context(), "role", claims["role"])
			ctx = context.WithValue(ctx, "user_id", claims["user_id"])

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
