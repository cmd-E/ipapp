package utils

import "strings"

var flags = []string{"-m"}

// FilterOutFlags filters out flags and it's values from provided arguments
func FilterOutFlags(providedIps []string) []string {
	for i := 0; i < len(providedIps); i++ {
		if isInTheFlagsList(providedIps[i]) {
			providedIps = append(providedIps[:i], providedIps[i+1:]...)
		}
	}
	return providedIps
}

func isInTheFlagsList(argToCheck string) bool {
	for _, v := range flags {
		if strings.HasPrefix(argToCheck, v) {
			return true
		}
	}
	return false
}
