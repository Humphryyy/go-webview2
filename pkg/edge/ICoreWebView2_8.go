//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2_8Vtbl struct {
	iCoreWebView2_7Vtbl
	AddIsMutedChanged ComProc
RemoveIsMutedChanged ComProc
GetIsMuted ComProc
PutIsMuted ComProc
AddIsDocumentPlayingAudioChanged ComProc
RemoveIsDocumentPlayingAudioChanged ComProc
GetIsDocumentPlayingAudio ComProc

}

type ICoreWebView2_8 struct {
	vtbl *iCoreWebView2_8Vtbl
}

func (i *ICoreWebView2_8) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_8) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_8() *ICoreWebView2_8 {
	var result *ICoreWebView2_8

	iidICoreWebView2_8 := NewGUID("")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_8)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (e *Chromium) GetICoreWebView2_8() *ICoreWebView2_8 {
	return e.webview.GetICoreWebView2_8()
}
