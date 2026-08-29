package middleware

import "net/http"

func GetUserID(r *http.Request) string {
	id, _ := r.Context().Value(UserIDContextKey).(string)
	return id
}

func GetToken(r *http.Request) string {
	t, _ := r.Context().Value(TokenContextKey).(string)
	return t
}

func GetRole(r *http.Request) string {
	role, _ := r.Context().Value(RoleContextKey).(string)
	return role
}