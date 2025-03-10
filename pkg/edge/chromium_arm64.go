//go:build windows
// +build windows

package edge

import (
	"unsafe"

	"github.com/Humphryyy/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

func (e *Chromium) SetSize(bounds w32.Rect) {
	if e.controller == nil {
		return
	}

	words := (*[2]uintptr)(unsafe.Pointer(&bounds))
	_, _, err := e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		words[0],
		words[1],
	)
	if err != windows.ERROR_SUCCESS {
		e.errorCallback(err)
	}
}
