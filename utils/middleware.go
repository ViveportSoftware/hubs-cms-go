package utils

import "strings"

func GetBearerToken(header string) (string, bool) {
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) != 2 {
		// Error: Bearer token not in proper format
		return "", false
	}

	result := strings.TrimSpace(splitToken[1])
	if len(result) == 0 {
		return "", false
	}

	return result, true
}
