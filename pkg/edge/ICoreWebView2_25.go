//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_25Vtbl struct {
	iCoreWebView2_24Vtbl
	AddSaveAsUIShowing ComProc
RemoveSaveAsUIShowing ComProc
ShowSaveAsUI ComProc

}

type ICoreWebView2_25 struct {
	vtbl *iCoreWebView2_25Vtbl
}

func (i *ICoreWebView2_25) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_25) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_25() *ICoreWebView2_25 {
	var result *ICoreWebView2_25

	iidICoreWebView2_25 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_25)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_25() *ICoreWebView2_25 {
	return e.webview.GetICoreWebView2_25()
}
