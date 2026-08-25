// Package authmw: middleware ตรวจ Bearer token — auth-specific ของ user-service.
package authmw

import (
	"context"
	"net/http"
	"strings"

	"github.com/kidkon/ecommerce/common/response"
	"github.com/kidkon/ecommerce/user-service/internal/token"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if !strings.HasPrefix(authz, "Bearer ") {
				response.Fail(w, response.Unauthorized("missing bearer token"))
				return
			}
			claims, err := token.Parse(strings.TrimPrefix(authz, "Bearer "), secret)
			if err != nil {
				response.Fail(w, response.Unauthorized("invalid or expired token"))
				return
			}
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFrom(ctx context.Context) *token.Claims {
	c, _ := ctx.Value(claimsKey).(*token.Claims)
	return c
}
