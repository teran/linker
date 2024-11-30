package models

import "time"

type Request struct {
	Timestamp  time.Time
	LinkID     string
	ClientIP   string
	CookieID   string
	UserAgent  string
	Parameters Parameters
	Referrer   string
}
