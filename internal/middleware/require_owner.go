package middleware

import "net/http"

func RequireOwner(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetRole(r) != "business_owner" {
			http.Redirect(w, r, "/access-denied", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}