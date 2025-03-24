//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_5Vtbl struct {
	iCoreWebView2_4Vtbl
	AddClientCertificateRequested ComProc
RemoveClientCertificateRequested ComProc

}

type ICoreWebView2_5 struct {
	vtbl *iCoreWebView2_5Vtbl
}

func (i *ICoreWebView2_5) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_5) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_5() *ICoreWebView2_5 {
	var result *ICoreWebView2_5

	iidICoreWebView2_5 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_5)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_5() *ICoreWebView2_5 {
	return e.webview.GetICoreWebView2_5()
}
