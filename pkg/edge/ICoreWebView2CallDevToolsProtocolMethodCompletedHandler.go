//go:build windows

package edge

import (
	"unsafe"
)

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandler struct {
	vtbl *ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl
	impl ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl
}

func (i *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler, errorCode uintptr, result *uint16) uintptr {
	// Convert result to string
	jsonResult := UTF16PtrToString(result)

	return this.impl.CallDevToolsProtocolMethodCompleted(errorCode, jsonResult)
}

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl interface {
	IUnknownImpl
	CallDevToolsProtocolMethodCompleted(errorCode uintptr, result string) uintptr
}

var _ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn = ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke, "_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke"),
}

func NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(impl ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl) *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler {
	return &ICoreWebView2CallDevToolsProtocolMethodCompletedHandler{
		vtbl: &_ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn,
		impl: impl,
	}
}
