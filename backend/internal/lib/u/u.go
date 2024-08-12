package u

import "strings"

func ParamFromPath(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
