//go:build windows

package edge

import (
	"unsafe"
)

type _ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2WebResourceResponseReceivedEventHandler struct {
	vtbl *_ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl
	impl _ICoreWebView2WebResourceResponseReceivedEventHandlerImpl
}

func (i *iCoreWebView2WebResourceResponseReceivedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func _ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2WebResourceResponseReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2WebResourceResponseReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownRelease(this *iCoreWebView2WebResourceResponseReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebResourceResponseReceivedEventHandlerInvoke(this *iCoreWebView2WebResourceResponseReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebResourceResponseReceivedEventArgs) uintptr {
	return this.impl.WebResourceResponseReceived(sender, args)
}

type _ICoreWebView2WebResourceResponseReceivedEventHandlerImpl interface {
	_IUnknownImpl
	WebResourceResponseReceived(sender *ICoreWebView2, args *ICoreWebView2WebResourceResponseReceivedEventArgs) uintptr
}

var ICoreWebView2WebResourceResponseReceivedEventHandlerFn = _ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2WebResourceResponseReceivedEventHandlerInvoke),
}

func newICoreWebView2WebResourceResponseReceivedEventHandler(impl _ICoreWebView2WebResourceResponseReceivedEventHandlerImpl) *iCoreWebView2WebResourceResponseReceivedEventHandler {
	return &iCoreWebView2WebResourceResponseReceivedEventHandler{
		vtbl: &ICoreWebView2WebResourceResponseReceivedEventHandlerFn,
		impl: impl,
	}
}
