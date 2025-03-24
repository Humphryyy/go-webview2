//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2ExecuteScriptResultVtbl struct {
	IUnknownVtbl
	GetSucceeded         ComProc
	GetResultAsJson      ComProc
	TryGetResultAsString ComProc
	GetException         ComProc
}

type ICoreWebView2ExecuteScriptResult struct {
	vtbl *ICoreWebView2ExecuteScriptResultVtbl
}

func (i *ICoreWebView2ExecuteScriptResult) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ExecuteScriptResult) GetSucceeded() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	_, _, err := i.vtbl.GetSucceeded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if err != (windows.ERROR_SUCCESS) {
		return false, err
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2ExecuteScriptResult) GetResultAsJson() (string, error) {
	// Create *uint16 to hold result
	var _jsonResult *uint16

	_, _, err := i.vtbl.GetResultAsJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_jsonResult)),
	)
	if err != (windows.ERROR_SUCCESS) {
		return "", err
	}
	// Get result and cleanup
	jsonResult := UTF16PtrToString(_jsonResult)
	CoTaskMemFree(unsafe.Pointer(_jsonResult))
	return jsonResult, nil
}

func (i *ICoreWebView2ExecuteScriptResult) TryGetResultAsString() (string, bool, error) {
	// Create *uint16 to hold result
	var _stringResult *uint16
	// Create int32 to hold bool result
	var _value int32

	_, _, err := i.vtbl.TryGetResultAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_stringResult)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if err != (windows.ERROR_SUCCESS) {
		return "", false, err
	}
	// Get result and cleanup
	stringResult := UTF16PtrToString(_stringResult)
	CoTaskMemFree(unsafe.Pointer(_stringResult))
	// Get result and cleanup
	value := _value != 0
	return stringResult, value, nil
}

func (i *ICoreWebView2ExecuteScriptResult) GetException() (*ICoreWebView2ScriptException, error) {

	var exception *ICoreWebView2ScriptException

	_, _, err := i.vtbl.GetException.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&exception)),
	)
	if err != (windows.ERROR_SUCCESS) {
		return nil, err
	}
	return exception, nil
}
