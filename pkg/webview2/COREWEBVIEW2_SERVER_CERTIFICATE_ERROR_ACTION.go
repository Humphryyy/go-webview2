//go:build windows

package webview2

type COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION uint32

const (
	COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION_ALWAYS_ALLOW = 0
	COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION_CANCEL       = 1
	COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION_DEFAULT      = 2
)
