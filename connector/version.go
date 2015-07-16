package main
import (
    "fmt"
    "runtime"
)

const BINARY_VERSION = "v0.5-26-g1d254a1"

func Version(app string) string {
    return fmt.Sprintf("%s %s (built w/%s)", app, BINARY_VERSION, runtime.Version())
}

