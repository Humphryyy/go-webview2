//go:build windows

package edge

import (
	"unsafe"
)

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandler struct {
	vtbl *ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl
	impl ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl
}

func (i *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func _ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownAddRef(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownRelease(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerInvoke(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, errorCode uintptr, result *IStream) uintptr {
	return this.impl.WebResourceResponseViewGetContentCompleted(errorCode, result)
}

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl interface {
	IUnknownImpl
	WebResourceResponseViewGetContentCompleted(errorCode uintptr, result *IStream) uintptr
}

var ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerFn = ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(_ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerInvoke),
}

func NewICoreWebView2WebResourceResponseViewGetContentCompletedHandler(impl ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl) *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler {
	return &ICoreWebView2WebResourceResponseViewGetContentCompletedHandler{
		vtbl: &ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerFn,
		impl: impl,
	}
}
