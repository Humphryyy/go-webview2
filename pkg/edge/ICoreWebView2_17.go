//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_17Vtbl struct {
	iCoreWebView2_16Vtbl
	PostSharedBufferToScript ComProc

}

type ICoreWebView2_17 struct {
	vtbl *iCoreWebView2_17Vtbl
}

func (i *ICoreWebView2_17) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_17) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_17() *ICoreWebView2_17 {
	var result *ICoreWebView2_17

	iidICoreWebView2_17 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_17)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_17() *ICoreWebView2_17 {
	return e.webview.GetICoreWebView2_17()
}
