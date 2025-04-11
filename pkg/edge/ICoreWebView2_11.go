//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type iCoreWebView2_11Vtbl struct {
	iCoreWebView2_10Vtbl
	CallDevToolsProtocolMethodForSession ComProc
	AddContextMenuRequested              ComProc
	RemoveContextMenuRequested           ComProc
}

type ICoreWebView2_11 struct {
	vtbl *iCoreWebView2_11Vtbl
}

func (i *ICoreWebView2_11) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2_11) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2) GetICoreWebView2_11() *ICoreWebView2_11 {
	var result *ICoreWebView2_11

	iidICoreWebView2_11 := NewGUID("{0be78e56-c193-4051-b943-23b460c08bdb}")
	_, _, _ = i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_11)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_11) CallDevToolsProtocolMethodForSession(sessionId string, methodName string, parametersAsJson string, handler *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) error {

	// Convert string 'sessionId' to *uint16
	_sessionId, err := UTF16PtrFromString(sessionId)
	if err != nil {
		return err
	}
	// Convert string 'methodName' to *uint16
	_methodName, err := UTF16PtrFromString(methodName)
	if err != nil {
		return err
	}
	// Convert string 'parametersAsJson' to *uint16
	_parametersAsJson, err := UTF16PtrFromString(parametersAsJson)
	if err != nil {
		return err
	}

	_, _, err = i.vtbl.CallDevToolsProtocolMethodForSession.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_sessionId)),
		uintptr(unsafe.Pointer(_methodName)),
		uintptr(unsafe.Pointer(_parametersAsJson)),
		uintptr(unsafe.Pointer(handler)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (e *Chromium) GetICoreWebView2_11() *ICoreWebView2_11 {
	return e.webview.GetICoreWebView2_11()
}
