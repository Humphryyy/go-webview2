//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_14Vtbl struct {
	iCoreWebView2_13Vtbl
	AddServerCertificateErrorDetected ComProc
RemoveServerCertificateErrorDetected ComProc
ClearServerCertificateErrorActions ComProc

}

type ICoreWebView2_14 struct {
	vtbl *iCoreWebView2_14Vtbl
}

func (i *ICoreWebView2_14) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_14) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_14() *ICoreWebView2_14 {
	var result *ICoreWebView2_14

	iidICoreWebView2_14 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_14)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_14() *ICoreWebView2_14 {
	return e.webview.GetICoreWebView2_14()
}
