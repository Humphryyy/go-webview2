//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_18Vtbl struct {
	iCoreWebView2_17Vtbl
	AddLaunchingExternalUriScheme ComProc
RemoveLaunchingExternalUriScheme ComProc

}

type ICoreWebView2_18 struct {
	vtbl *iCoreWebView2_18Vtbl
}

func (i *ICoreWebView2_18) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_18) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_18() *ICoreWebView2_18 {
	var result *ICoreWebView2_18

	iidICoreWebView2_18 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_18)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_18() *ICoreWebView2_18 {
	return e.webview.GetICoreWebView2_18()
}
