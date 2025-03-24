//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_12Vtbl struct {
	iCoreWebView2_11Vtbl
	AddStatusBarTextChanged ComProc
RemoveStatusBarTextChanged ComProc
GetStatusBarText ComProc

}

type ICoreWebView2_12 struct {
	vtbl *iCoreWebView2_12Vtbl
}

func (i *ICoreWebView2_12) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_12) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_12() *ICoreWebView2_12 {
	var result *ICoreWebView2_12

	iidICoreWebView2_12 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_12)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_12() *ICoreWebView2_12 {
	return e.webview.GetICoreWebView2_12()
}
