//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_24Vtbl struct {
	iCoreWebView2_23Vtbl
	AddNotificationReceived ComProc
RemoveNotificationReceived ComProc

}

type ICoreWebView2_24 struct {
	vtbl *iCoreWebView2_24Vtbl
}

func (i *ICoreWebView2_24) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_24) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_24() *ICoreWebView2_24 {
	var result *ICoreWebView2_24

	iidICoreWebView2_24 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_24)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_24() *ICoreWebView2_24 {
	return e.webview.GetICoreWebView2_24()
}
