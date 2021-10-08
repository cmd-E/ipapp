package printer

import "fmt"

const format = "%-29s %v\n"

// PrintFormatted prints provided values in special format
func PrintFormatted(definition string, value interface{}) {
	fmt.Printf(format, definition, value)
}
