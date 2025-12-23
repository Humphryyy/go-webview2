package edge

import (
	"unsafe"
)

type ICoreWebView2ExecuteScriptCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ExecuteScriptCompletedHandler struct {
	vtbl *ICoreWebView2ExecuteScriptCompletedHandlerVtbl
	impl ICoreWebView2ExecuteScriptCompletedHandlerImpl
}

func (i *ICoreWebView2ExecuteScriptCompletedHandler) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2ExecuteScriptCompletedHandler) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ExecuteScriptCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef(this *ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease(this *ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2ExecuteScriptCompletedHandlerInvoke(this *ICoreWebView2ExecuteScriptCompletedHandler, errorCode uintptr, result *uint16) uintptr {
	jsonResult := UTF16PtrToString(result)

	return this.impl.ExecuteScriptCompleted(errorCode, jsonResult)
}

type ICoreWebView2ExecuteScriptCompletedHandlerImpl interface {
	IUnknownImpl
	ExecuteScriptCompleted(errorCode uintptr, executedScript string) uintptr
}

var _ICoreWebView2ExecuteScriptCompletedHandlerFn = ICoreWebView2ExecuteScriptCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerInvoke),
}

func NewICoreWebView2ExecuteScriptCompletedHandler(impl ICoreWebView2ExecuteScriptCompletedHandlerImpl) *ICoreWebView2ExecuteScriptCompletedHandler {
	return &ICoreWebView2ExecuteScriptCompletedHandler{
		vtbl: &_ICoreWebView2ExecuteScriptCompletedHandlerFn,
		impl: impl,
	}
}
