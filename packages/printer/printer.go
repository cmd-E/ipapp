package printer

import "fmt"

const format = "%-29s %v\n"

func PrintFormatted(definition string, value interface{}) {
	fmt.Printf(format, definition, value)
}
