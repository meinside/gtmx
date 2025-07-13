// helpers.go

package main

import "strings"

// escape string before string interpolation
func escape(str string) string {
	return strings.ReplaceAll(
		str,
		`%`,
		`%%`,
	)
}
