// helpers.go
package helpers

func PadString(s string, width int) string {
	if len(s) < width {
		return s + spaces(width-len(s))
	}
	return s
}

func spaces(n int) string {
	space := ""
	for i := 0; i < n; i++ {
		space += " "
	}
	return space
}
