package fb

import (
	"encoding/base64"
)

// We need a URL-safe encoding, but '-' has a special meaning in FB doc IDs, so
// we use '=' instead (which is available since we don't use padding)
const b64alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789=_"

var b64encoder = base64.NewEncoding(b64alphabet).WithPadding(base64.NoPadding)
