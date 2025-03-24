//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2ScriptExceptionVtbl struct {
	IUnknownVtbl
	GetLineNumber   ComProc
	GetColumnNumber ComProc
	GetName         ComProc
	GetMessage      ComProc
	GetToJson       ComProc
}

type ICoreWebView2ScriptException struct {
	vtbl *ICoreWebView2ScriptExceptionVtbl
}

func (i *ICoreWebView2ScriptException) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ScriptException) GetLineNumber() (uint32, error) {

	var value uint32

	hr, _, err := i.vtbl.GetLineNumber.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ScriptException) GetColumnNumber() (uint32, error) {

	var value uint32

	hr, _, err := i.vtbl.GetColumnNumber.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ScriptException) GetName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, err := i.vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2ScriptException) GetMessage() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, err := i.vtbl.GetMessage.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2ScriptException) GetToJson() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, err := i.vtbl.GetToJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}
