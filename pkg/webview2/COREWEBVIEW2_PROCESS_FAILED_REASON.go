//go:build windows

package webview2

type COREWEBVIEW2_PROCESS_FAILED_REASON uint32

const (
	COREWEBVIEW2_PROCESS_FAILED_REASON_UNEXPECTED      = 0
	COREWEBVIEW2_PROCESS_FAILED_REASON_UNRESPONSIVE    = 1
	COREWEBVIEW2_PROCESS_FAILED_REASON_TERMINATED      = 2
	COREWEBVIEW2_PROCESS_FAILED_REASON_CRASHED         = 3
	COREWEBVIEW2_PROCESS_FAILED_REASON_LAUNCH_FAILED   = 4
	COREWEBVIEW2_PROCESS_FAILED_REASON_OUT_OF_MEMORY   = 5
	COREWEBVIEW2_PROCESS_FAILED_REASON_PROFILE_DELETED = 6
)
