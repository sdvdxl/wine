package constant

import (
	"time"
)

const (
	LOGIN_ERROR_MSG               = "_LOGIN_ERROR_MSG"
	LOGIN_PID                     = "_LOGIN_PID"
	LOGIN_USER                    = "_LOGIN_USER"
	LOGIN_TIMESTAMP               = "_LOGIN_TIMESTAMP"
	LOGIN_EXPIRED_TIME            = time.Minute * 60
	LOGIN_CAPTCHA                 = "LOGIN_CAPTCHA"
	PHONE_CAPTCHA_SECONDS         = 60
	PHONE_CAPTCHA_LAST            = "_PHONE_CAPTCHA_LAST"
	PHONE_CAPTCHA                 = "_PHONE_CAPTCHA"
	PHONE_CAPTCHA_EXPIRED_MINUTES = 5
	PHONE_NUMBER                  = "_PHNOE_NUMBER"
)
