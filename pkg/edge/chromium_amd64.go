//go:build windows
// +build windows

package edge

import (
	"errors"
	"unsafe"

	"github.com/Humphryyy/go-webview2/internal/w32"
	"golang.org/x/sys/windows"
)

func (e *Chromium) SetSize(bounds w32.Rect) {
	if e.controller == nil {
		return
	}

	_, _, err := e.controller.vtbl.PutBounds.Call(
		uintptr(unsafe.Pointer(e.controller)),
		uintptr(unsafe.Pointer(&bounds)),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		e.errorCallback(err)
	}
}
