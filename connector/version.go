package main
import (
    "fmt"
    "runtime"
)

const BINARY_VERSION = "v0.5-28-g605fe1d"

func Version(app string) string {
    return fmt.Sprintf("%s %s (built w/%s)", app, BINARY_VERSION, runtime.Version())
}

