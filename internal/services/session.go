package services

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Store is the session store
var Store = sessions.NewCookieStore([]byte("super-secret-key-change-this"))

func Init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1800, // 30 mins
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

// GetSession returns the current session
func GetSession(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, "od_session")
	return session
}
