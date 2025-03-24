//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_16Vtbl struct {
	iCoreWebView2_15Vtbl
	Print ComProc
ShowPrintUI ComProc
PrintToPdfStream ComProc

}

type ICoreWebView2_16 struct {
	vtbl *iCoreWebView2_16Vtbl
}

func (i *ICoreWebView2_16) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_16) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_16() *ICoreWebView2_16 {
	var result *ICoreWebView2_16

	iidICoreWebView2_16 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_16)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_16() *ICoreWebView2_16 {
	return e.webview.GetICoreWebView2_16()
}
