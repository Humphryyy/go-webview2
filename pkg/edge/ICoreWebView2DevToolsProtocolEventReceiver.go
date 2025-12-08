//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2DevToolsProtocolEventReceiverVtbl struct {
	_IUnknownVtbl
	AddDevToolsProtocolEventReceived    ComProc
	RemoveDevToolsProtocolEventReceived ComProc
}

type ICoreWebView2DevToolsProtocolEventReceiver struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceiverVtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceiver) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

type DevToolsProtocolEventReceivedHandler struct {
	resultFunc func(sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr
}

func (d *DevToolsProtocolEventReceivedHandler) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (d *DevToolsProtocolEventReceivedHandler) AddRef() uintptr {
	return 1
}

func (d *DevToolsProtocolEventReceivedHandler) Release() uintptr {
	return 1
}

func (d *DevToolsProtocolEventReceivedHandler) DevToolsProtocolEventReceived(sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr {
	return d.resultFunc(sender, args)
}

var receiverPool = make([]*iCoreWebView2DevToolsProtocolEventReceivedEventHandler, 0)

func (i *ICoreWebView2DevToolsProtocolEventReceiver) AddDevToolsProtocolEventReceived(handler func(sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr) (EventRegistrationToken, error) {
	handlerImpl := &DevToolsProtocolEventReceivedHandler{resultFunc: handler}
	eventHandler := NewICoreWebView2DevToolsProtocolEventReceivedEventHandler(handlerImpl)

	receiverPool = append(receiverPool, eventHandler)

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddDevToolsProtocolEventReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2DevToolsProtocolEventReceiver) RemoveDevToolsProtocolEventReceived(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveDevToolsProtocolEventReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
