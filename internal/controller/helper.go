package controller

import (
	"strings"
)

func isAllowedEtcdKey(key string) bool {
	key = strings.TrimSpace(key)

	if !strings.HasPrefix(key, "/system/") {
		return false
	}

	// extract next segment after /system/
	segments := strings.SplitN(strings.TrimPrefix(key, "/system/"), "/", 2)
	if len(segments) == 0 {
		return false
	}
	serviceName := segments[0]
	return strings.HasSuffix(serviceName, "-services")
}
