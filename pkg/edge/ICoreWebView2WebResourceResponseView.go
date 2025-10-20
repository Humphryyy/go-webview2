//go:build windows

package edge

import (
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2WebResourceResponseViewVtbl struct {
	_IUnknownVtbl
	GetHeaders      ComProc
	GetStatusCode   ComProc
	GetReasonPhrase ComProc
	GetContent      ComProc
}

type ICoreWebView2WebResourceResponseView struct {
	vtbl *ICoreWebView2WebResourceResponseViewVtbl
}

func (i *ICoreWebView2WebResourceResponseView) AddRef() uintptr {
	refCounter, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2WebResourceResponseView) Release() uintptr {
	refCounter, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2WebResourceResponseView) GetHeaders() (*ICoreWebView2HttpResponseHeaders, error) {

	var headers *ICoreWebView2HttpResponseHeaders

	hr, _, _ := i.vtbl.GetHeaders.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return headers, nil
}

func (i *ICoreWebView2WebResourceResponseView) GetStatusCode() (int, error) {

	var statusCode int

	hr, _, _ := i.vtbl.GetStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(statusCode),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return statusCode, nil
}

func (i *ICoreWebView2WebResourceResponseView) GetReasonPhrase() (string, error) {
	// Create *uint16 to hold result
	var _reasonPhrase *uint16

	hr, _, _ := i.vtbl.GetReasonPhrase.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_reasonPhrase)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	reasonPhrase := UTF16PtrToString(_reasonPhrase)
	CoTaskMemFree(unsafe.Pointer(_reasonPhrase))
	return reasonPhrase, nil
}

type WebResourceResponseViewGetContentCompletedHandler struct {
	resultFunc func(errorCode uintptr, result *IStream) uintptr
}

func (*WebResourceResponseViewGetContentCompletedHandler) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (*WebResourceResponseViewGetContentCompletedHandler) AddRef() uintptr {
	return 1
}

func (*WebResourceResponseViewGetContentCompletedHandler) Release() uintptr {
	return 1
}

func (h *WebResourceResponseViewGetContentCompletedHandler) WebResourceResponseViewGetContentCompleted(errorCode uintptr, result *IStream) uintptr {
	return h.resultFunc(errorCode, result)
}

// Keep a pool of handlers to prevent garbage collection
var contentHandlerPool = make([]*ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, 0, 10000)
var contentHandlerPoolMutex sync.Mutex

func (i *ICoreWebView2WebResourceResponseView) GetContent(callback func(errorCode uintptr, result *IStream) uintptr) error {
	handlerImpl := &WebResourceResponseViewGetContentCompletedHandler{}

	handlerImpl.resultFunc = callback

	handler := NewICoreWebView2WebResourceResponseViewGetContentCompletedHandler(handlerImpl)

	contentHandlerPoolMutex.Lock()
	contentHandlerPool = append(contentHandlerPool, handler)

	if len(contentHandlerPool) > 10000 {
		half := len(contentHandlerPool) / 2
		tmp := make([]*ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, len(contentHandlerPool)-half)
		copy(tmp, contentHandlerPool[half:])
		contentHandlerPool = tmp
	}

	contentHandlerPoolMutex.Unlock()

	hr, _, _ := i.vtbl.GetContent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}

	return nil
}
