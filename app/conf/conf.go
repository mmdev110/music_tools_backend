package conf

import (
	"time"
)

var TokenDuration = 1 * time.Minute

// var RefreshDuration = 1 * 24 * time.Hour * 30
var RefreshDuration = 10 * time.Minute
var CookieDomain = "localhost"
var SessionID_KEY = "session_id"
