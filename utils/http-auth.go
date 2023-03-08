package utils

import (
	"encoding/base64"
	"fmt"
)

func ToBasicAuth(role, password string) string {
	s := fmt.Sprint(role, ":", password)
	b := base64.RawStdEncoding.EncodeToString([]byte(s))
	return fmt.Sprint("Basic ", b)
}
