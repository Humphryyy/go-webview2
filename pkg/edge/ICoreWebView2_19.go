//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_19Vtbl struct {
	iCoreWebView2_18Vtbl
	GetMemoryUsageTargetLevel ComProc
PutMemoryUsageTargetLevel ComProc

}

type ICoreWebView2_19 struct {
	vtbl *iCoreWebView2_19Vtbl
}

func (i *ICoreWebView2_19) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_19) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_19() *ICoreWebView2_19 {
	var result *ICoreWebView2_19

	iidICoreWebView2_19 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_19)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_19() *ICoreWebView2_19 {
	return e.webview.GetICoreWebView2_19()
}
