package models

type Request struct {
	LinkID     string
	ClientIP   string
	CookieID   string
	UserAgent  string
	Parameters Parameters
}
