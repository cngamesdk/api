package cache

import (
	"fmt"
	"strings"
)

func BuildCacheKey(req ...interface{}) string {
	var build []string
	for _, item := range req {
		build = append(build, fmt.Sprintf("%v", item))
	}
	return strings.Join(build, "-")
}
