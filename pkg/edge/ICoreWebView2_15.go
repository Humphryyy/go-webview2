//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_15Vtbl struct {
	iCoreWebView2_14Vtbl
	AddFaviconChanged ComProc
RemoveFaviconChanged ComProc
GetFaviconUri ComProc
GetFavicon ComProc

}

type ICoreWebView2_15 struct {
	vtbl *iCoreWebView2_15Vtbl
}

func (i *ICoreWebView2_15) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_15) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_15() *ICoreWebView2_15 {
	var result *ICoreWebView2_15

	iidICoreWebView2_15 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_15)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_15() *ICoreWebView2_15 {
	return e.webview.GetICoreWebView2_15()
}
