//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_23Vtbl struct {
	iCoreWebView2_22Vtbl
	PostWebMessageAsJsonWithAdditionalObjects ComProc

}

type ICoreWebView2_23 struct {
	vtbl *iCoreWebView2_23Vtbl
}

func (i *ICoreWebView2_23) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_23) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_23() *ICoreWebView2_23 {
	var result *ICoreWebView2_23

	iidICoreWebView2_23 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_23)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_23() *ICoreWebView2_23 {
	return e.webview.GetICoreWebView2_23()
}
