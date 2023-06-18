package helpers

import "os"

func Env(key string) string {
	strVal, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	return strVal
}
