// +build !linux,!darwin

package utils

import (
    "fmt"
)

func dup2(oldfd int, newfd int) error {
    return fmt.Errorf("dup2 is not supported on this platform")
}
