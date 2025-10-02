//go:build windows

package edge

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2_2Vtbl struct {
	QueryInterface ComProc
	AddRef         ComProc
	Release        ComProc
	// ICoreWebView2 methods
	GetSettings                            ComProc
	GetSource                              ComProc
	Navigate                               ComProc
	NavigateToString                       ComProc
	AddNavigationStarting                  ComProc
	RemoveNavigationStarting               ComProc
	AddContentLoading                      ComProc
	RemoveContentLoading                   ComProc
	AddSourceChanged                       ComProc
	RemoveSourceChanged                    ComProc
	AddHistoryChanged                      ComProc
	RemoveHistoryChanged                   ComProc
	AddNavigationCompleted                 ComProc
	RemoveNavigationCompleted              ComProc
	AddFrameNavigationStarting             ComProc
	RemoveFrameNavigationStarting          ComProc
	AddFrameNavigationCompleted            ComProc
	RemoveFrameNavigationCompleted         ComProc
	AddScriptDialogOpening                 ComProc
	RemoveScriptDialogOpening              ComProc
	AddPermissionRequested                 ComProc
	RemovePermissionRequested              ComProc
	AddProcessFailed                       ComProc
	RemoveProcessFailed                    ComProc
	AddScriptToExecuteOnDocumentCreated    ComProc
	RemoveScriptToExecuteOnDocumentCreated ComProc
	ExecuteScript                          ComProc
	CapturePreview                         ComProc
	Reload                                 ComProc
	PostWebMessageAsJSON                   ComProc
	PostWebMessageAsString                 ComProc
	AddWebMessageReceived                  ComProc
	RemoveWebMessageReceived               ComProc
	CallDevToolsProtocolMethod             ComProc
	GetBrowserProcessID                    ComProc
	GetCanGoBack                           ComProc
	GetCanGoForward                        ComProc
	GoBack                                 ComProc
	GoForward                              ComProc
	GetDevToolsProtocolEventReceiver       ComProc
	Stop                                   ComProc
	AddNewWindowRequested                  ComProc
	RemoveNewWindowRequested               ComProc
	AddDocumentTitleChanged                ComProc
	RemoveDocumentTitleChanged             ComProc
	GetDocumentTitle                       ComProc
	AddHostObjectToScript                  ComProc
	RemoveHostObjectFromScript             ComProc
	OpenDevToolsWindow                     ComProc
	AddContainsFullScreenElementChanged    ComProc
	RemoveContainsFullScreenElementChanged ComProc
	GetContainsFullScreenElement           ComProc
	AddWebResourceRequested                ComProc
	RemoveWebResourceRequested             ComProc
	AddWebResourceRequestedFilter          ComProc
	RemoveWebResourceRequestedFilter       ComProc
	AddWindowCloseRequested                ComProc
	RemoveWindowCloseRequested             ComProc
	// ICoreWebView2_2 methods
	AddWebResourceResponseReceived    ComProc
	RemoveWebResourceResponseReceived ComProc
	NavigateWithWebResourceRequest    ComProc
	AddDomContentLoaded               ComProc
	RemoveDomContentLoaded            ComProc
	GetCookieManager                  ComProc
	GetEnvironment                    ComProc
}

type ICoreWebView2_2 struct {
	vtbl *iCoreWebView2_2Vtbl
}

// AddRef increments the reference count of the ICoreWebView2_2 interface
func (i *ICoreWebView2_2) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

// Release decrements the reference count of the ICoreWebView2_2 interface
func (i *ICoreWebView2_2) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_2) GetCookieManager() (*ICoreWebView2CookieManager, error) {
	var cookieManager *ICoreWebView2CookieManager
	hr, _, _ := i.vtbl.GetCookieManager.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&cookieManager)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return cookieManager, nil
}

func (i *ICoreWebView2) GetICoreWebView2_2() *ICoreWebView2_2 {
	var result *ICoreWebView2_2

	iidICoreWebView2_2 := NewGUID("{9E8F0CF8-E670-4B5E-B2BC-73E061E3184C}")
	hr, _, _ := i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_2)),
		uintptr(unsafe.Pointer(&result)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil
	}
	return result
}

func (e *Chromium) GetICoreWebView2_2() *ICoreWebView2_2 {
	return e.webview.GetICoreWebView2_2()
}

func (i *ICoreWebView2) AddWebResourceResponseReceived(handler *iCoreWebView2WebResourceResponseReceivedEventHandler, token *_EventRegistrationToken) error {
	webview2 := i.GetICoreWebView2_2()
	if webview2 == nil {
		return fmt.Errorf("ICoreWebView2_2 interface not available")
	}

	hr, _, _ := webview2.vtbl.AddWebResourceResponseReceived.Call(
		uintptr(unsafe.Pointer(webview2)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}
