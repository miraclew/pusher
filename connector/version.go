package main
import (
    "fmt"
    "runtime"
)

const BINARY_VERSION = "v0.5-44-gc34d741"

func Version(app string) string {
    return fmt.Sprintf("%s %s (built w/%s)", app, BINARY_VERSION, runtime.Version())
}

