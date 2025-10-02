//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2HttpResponseHeadersVtbl struct {
	_IUnknownVtbl
	AppendHeader ComProc
	Contains     ComProc
	GetHeader    ComProc
	GetHeaders   ComProc
	GetIterator  ComProc
}

type ICoreWebView2HttpResponseHeaders struct {
	vtbl *_ICoreWebView2HttpResponseHeadersVtbl
}

func (i *ICoreWebView2HttpResponseHeaders) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2HttpResponseHeaders) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2HttpResponseHeaders) AppendHeader(name string, value string) error {
	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.AppendHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

// GetHeader returns the value of the specified header. If the header is not found
// ERROR_ELEMENT_NOT_FOUND is returned as error.
func (i *ICoreWebView2HttpResponseHeaders) GetHeader(name string) (string, error) {
	_name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return "", err
	}

	var _value *uint16
	hr, _, _ := i.vtbl.GetHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}

	value := windows.UTF16PtrToString(_value)
	windows.CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

// GetIterator returns an iterator over the collection of response headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2HttpResponseHeaders) GetIterator() (*ICoreWebView2HttpHeadersCollectionIterator, error) {
	var headers *ICoreWebView2HttpHeadersCollectionIterator
	hr, _, _ := i.vtbl.GetIterator.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}

	return headers, nil
}
