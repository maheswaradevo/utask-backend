package helpers

import (
	"fmt"
	"strings"
)

func CacheWithPrefix(prefix, param string) string {
	param = HashSHA256(param)
	return fmt.Sprintf("%s:%s", strings.ToLower(prefix), param)
}
