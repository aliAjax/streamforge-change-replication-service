package sourceadapter

import "strings"

func RedactEndpoint(endpoint string) string {
	if i := strings.Index(endpoint, "@"); i >= 0 {
		return "***@" + endpoint[i+1:]
	}
	return endpoint
}
func AllowedColumn(name string, excluded []string) bool {
	for _, v := range excluded {
		if v == name {
			return false
		}
	}
	return true
}
