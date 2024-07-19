// +build linux darwin
package utils

import (
    // "os"
    // "syscall"
	"golang.org/x/sys/unix"
)

func dup2(oldfd int, newfd int) error {
    return unix.Dup2(oldfd, newfd)
    // return syscall.Dup2(oldfd, newfd)
}