package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/retail-core/xsell-web/internal/session"
)

type contextKey string

const TokenContextKey contextKey = "token"
const UserIDContextKey contextKey = "user_id"
const RoleContextKey contextKey = "role"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := session.GetToken(r)
		if token == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, err := session.DecodeAndVerify(token, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			session.ClearToken(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), TokenContextKey, token)
		ctx = context.WithValue(ctx, UserIDContextKey, claims.Sub)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}