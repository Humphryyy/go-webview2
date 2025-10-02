//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2WebResourceResponseReceivedEventArgsVtbl struct {
	_IUnknownVtbl
	GetRequest  ComProc
	GetResponse ComProc
}

type ICoreWebView2WebResourceResponseReceivedEventArgs struct {
	vtbl *_ICoreWebView2WebResourceResponseReceivedEventArgsVtbl
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return refCounter
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest, error) {

	var request *ICoreWebView2WebResourceRequest

	hr, _, _ := i.vtbl.GetRequest.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&request)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return request, nil
}

func (i *ICoreWebView2WebResourceResponseReceivedEventArgs) GetResponse() (*ICoreWebView2WebResourceResponseView, error) {

	var response *ICoreWebView2WebResourceResponseView

	hr, _, _ := i.vtbl.GetResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return response, nil
}
