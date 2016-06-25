package context

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

const (
	TIOCGWINSZ = 1074295912
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TerminalWidth() int {
	sizeobj, err := GetWinsize()
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return int(sizeobj.Col)
}

func GetWinsize() (*winsize, error) {
	ws := new(winsize)

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		return nil, errors.New(errno.Error())
	}
	return ws, nil
}
