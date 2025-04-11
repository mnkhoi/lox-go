package lox

import (
	"fmt"
)

func error(line int, msg string) {
	report(line, "", msg)
}

func report(line int, where, msg string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, msg)
}


