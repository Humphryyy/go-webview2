//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PermissionSettingCollectionViewVtbl struct {
	IUnknownVtbl
	GetValueAtIndex ComProc
	GetCount        ComProc
}

type ICoreWebView2PermissionSettingCollectionView struct {
	Vtbl *ICoreWebView2PermissionSettingCollectionViewVtbl
}

func (i *ICoreWebView2PermissionSettingCollectionView) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2PermissionSettingCollectionView) GetValueAtIndex(index uint32) (*ICoreWebView2PermissionSetting, error) {

	var permissionSetting *ICoreWebView2PermissionSetting

	hr, _, err := i.Vtbl.GetValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
		uintptr(unsafe.Pointer(&permissionSetting)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return permissionSetting, err
}

func (i *ICoreWebView2PermissionSettingCollectionView) GetCount() (*uint32, error) {

	var value *uint32

	hr, _, err := i.Vtbl.GetCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
