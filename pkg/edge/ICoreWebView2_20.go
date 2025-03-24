//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_20Vtbl struct {
	iCoreWebView2_19Vtbl
	GetFrameId ComProc

}

type ICoreWebView2_20 struct {
	vtbl *iCoreWebView2_20Vtbl
}

func (i *ICoreWebView2_20) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_20) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_20() *ICoreWebView2_20 {
	var result *ICoreWebView2_20

	iidICoreWebView2_20 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_20)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_20() *ICoreWebView2_20 {
	return e.webview.GetICoreWebView2_20()
}
