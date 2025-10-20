//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetParameterObjectAsJson ComProc
	GetSessionID             ComProc
}

type ICoreWebView2DevToolsProtocolEventReceivedEventArgs struct {
	Vtbl *iCoreWebView2DevToolsProtocolEventReceivedEventArgsVtbl
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) GetParameterObjectAsJson() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetParameterObjectAsJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}

	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) GetSessionID() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetSessionID.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}

	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}
