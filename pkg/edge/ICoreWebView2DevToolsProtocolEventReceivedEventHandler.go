//go:build windows

package edge

import (
	"unsafe"
)

type iCoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2DevToolsProtocolEventReceivedEventHandler struct {
	Vtbl *iCoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl
	impl ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl
}

func (i *iCoreWebView2DevToolsProtocolEventReceivedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2DevToolsProtocolEventReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2DevToolsProtocolEventReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownRelease(this *iCoreWebView2DevToolsProtocolEventReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerInvoke(this *iCoreWebView2DevToolsProtocolEventReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr {
	return this.impl.DevToolsProtocolEventReceived(sender, args)
}

type ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl interface {
	IUnknownImpl
	DevToolsProtocolEventReceived(sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr
}

var ICoreWebView2DevToolsProtocolEventReceivedEventHandlerFn = iCoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerInvoke),
}

func NewICoreWebView2DevToolsProtocolEventReceivedEventHandler(impl ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl) *iCoreWebView2DevToolsProtocolEventReceivedEventHandler {
	return &iCoreWebView2DevToolsProtocolEventReceivedEventHandler{
		Vtbl: &ICoreWebView2DevToolsProtocolEventReceivedEventHandlerFn,
		impl: impl,
	}
}
