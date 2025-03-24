//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_13Vtbl struct {
	iCoreWebView2_12Vtbl
	GetProfile ComProc

}

type ICoreWebView2_13 struct {
	vtbl *iCoreWebView2_13Vtbl
}

func (i *ICoreWebView2_13) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_13) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_13() *ICoreWebView2_13 {
	var result *ICoreWebView2_13

	iidICoreWebView2_13 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_13)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_13() *ICoreWebView2_13 {
	return e.webview.GetICoreWebView2_13()
}
