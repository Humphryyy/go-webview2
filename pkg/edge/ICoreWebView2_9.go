//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_9Vtbl struct {
	iCoreWebView2_8Vtbl
	AddIsDefaultDownloadDialogOpenChanged ComProc
RemoveIsDefaultDownloadDialogOpenChanged ComProc
GetIsDefaultDownloadDialogOpen ComProc
OpenDefaultDownloadDialog ComProc
CloseDefaultDownloadDialog ComProc
GetDefaultDownloadDialogCornerAlignment ComProc
PutDefaultDownloadDialogCornerAlignment ComProc
GetDefaultDownloadDialogMargin ComProc
PutDefaultDownloadDialogMargin ComProc

}

type ICoreWebView2_9 struct {
	vtbl *iCoreWebView2_9Vtbl
}

func (i *ICoreWebView2_9) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_9) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_9() *ICoreWebView2_9 {
	var result *ICoreWebView2_9

	iidICoreWebView2_9 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_9)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_9() *ICoreWebView2_9 {
	return e.webview.GetICoreWebView2_9()
}
