package session

import "net/http"

type Session struct {
	name     string
	httpOnly bool
	secure   bool
	maxAge   int
}

func New(
	name string,
	httpOnly bool,
	secure bool,
	maxAge int,
) *Session {
	return &Session{
		name:     name,
		httpOnly: httpOnly,
		secure:   secure,
		maxAge:   maxAge,
	}
}

func (s *Session) Create(w http.ResponseWriter, sessionId string) {
	cookie := &http.Cookie{
		Name:     s.name,
		Value:    sessionId,
		Path:     "/",
		HttpOnly: s.httpOnly,
		Secure:   s.secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   s.maxAge,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusCreated)
}
