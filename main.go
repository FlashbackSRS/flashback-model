package fb

import (
	"encoding/base64"
	"time"
)

var b64encoder *base64.Encoding = base64.URLEncoding.WithPadding(base64.NoPadding)

func TimesEqual(t1, t2 *time.Time) bool {
	if t1 == nil && t2 == nil {
		return true
	}
	if t1 == nil || t2 == nil {
		return false
	}
	return t2.Equal(*t2)
}
