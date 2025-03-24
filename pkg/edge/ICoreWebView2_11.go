//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_11Vtbl struct {
	iCoreWebView2_10Vtbl
	CallDevToolsProtocolMethodForSession ComProc
AddContextMenuRequested ComProc
RemoveContextMenuRequested ComProc

}

type ICoreWebView2_11 struct {
	vtbl *iCoreWebView2_11Vtbl
}

func (i *ICoreWebView2_11) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_11) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_11() *ICoreWebView2_11 {
	var result *ICoreWebView2_11

	iidICoreWebView2_11 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_11)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_11() *ICoreWebView2_11 {
	return e.webview.GetICoreWebView2_11()
}
