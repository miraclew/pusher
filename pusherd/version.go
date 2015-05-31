package main

import (
	"fmt"
	"runtime"
)

func Version(app string) string {
	return fmt.Sprintf("%s %s (built w/%s)", app, BINARY_VERSION, runtime.Version())
}
