package middleware

import (
	"fmt"
	"net/http"

	"github.com/prajwal-huggi/backend_go/internal/domain"
)

func RequireRole(allowedRoles ...domain.RoleType) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			roleStr, ok:= r.Context().Value("role").(string)
			if !ok {
				fmt.Printf(`"%v": Role not found in context`, roleStr)
				http.Error(w, "unauthorized", 401)// concept unauthenticated
				return
			}

			role:= domain.RoleType(roleStr)

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "forbidden", 403)// concept unauthorized
		})
	}
}