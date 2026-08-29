package session

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/retail-core/xsell-web/internal/clients"
)

const cookieName = "jarabiz_session"
const storeIDCookieName = "jarabiz_active_store"
const activeStoreCookieName = "jarabiz_active_store"

func SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true once you have HTTPS in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})
}

func GetToken(r *http.Request) string {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func ClearToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
	})
}

func SetActiveStore(w http.ResponseWriter, store clients.Store) error {
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     storeIDCookieName,
		Value:    store.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     activeStoreCookieName,
		Value:    base64.URLEncoding.EncodeToString(data), // cookies need safe encoding
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 7,
	})
	return nil
}

func GetActiveStore(r *http.Request) (*clients.Store, error) {
	cookie, err := r.Cookie(activeStoreCookieName)
	if err != nil {
		return nil, err
	}
	data, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}
	var store clients.Store
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return &store, nil
}

func GetActiveStoreID(r *http.Request) string {
	cookie, err := r.Cookie(storeIDCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}