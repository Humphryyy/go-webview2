//go:build windows && !native_webview2loader

package webviewloader

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/Humphryyy/go-webview2/pkg/combridge"
	"golang.org/x/sys/windows"
)

func init() {
	UsingGoWebview2Loader = true
	preventEnvAndRegistryOverrides()
}

type webView2RunTimeType int32

const (
	webView2RunTimeTypeInstalled       webView2RunTimeType = 0x00
	webView2RunTimeTypeRedistributable webView2RunTimeType = 0x01
)

// CreateCoreWebView2Environment creates an evergreen WebView2 Environment using the installed WebView2 Runtime version.
//
// This is equivalent to running CreateCoreWebView2EnvironmentWithOptions without any options.
// For more information, see CreateCoreWebView2EnvironmentWithOptions.
func CreateCoreWebView2Environment(environmentCompletedHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) error {
	return CreateCoreWebView2EnvironmentWithOptions(environmentCompletedHandler)
}

// CreateCoreWebView2EnvironmentWithOptions creates an environment with a custom version of WebView2 Runtime,
// user data folder, and with or without additional options.
//
// See https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/webview2-idl?#createcorewebview2environmentwithoptions
func CreateCoreWebView2EnvironmentWithOptions(environmentCompletedHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, opts ...option) error {
	var params environmentOptions
	for _, opt := range opts {
		opt(&params)
	}

	var err error
	var dllPath string
	var runtimeType webView2RunTimeType
	if browserExecutableFolder := params.browserExecutableFolder; browserExecutableFolder != "" {
		runtimeType = webView2RunTimeTypeRedistributable
		dllPath, err = findEmbeddedClientDll(browserExecutableFolder)
	} else {
		runtimeType = webView2RunTimeTypeInstalled
		dllPath, _, err = findInstalledClientDll(params.preferCanary)
	}

	if err != nil {
		return err
	}

	return createWebViewEnvironmentWithClientDll(dllPath, runtimeType, params.userDataFolder,
		&params, environmentCompletedHandler)
}

func createWebViewEnvironmentWithClientDll(lpLibFileName string, runtimeType webView2RunTimeType, userDataFolder string,
	envOptions *environmentOptions, envCompletedHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) error {

	if !filepath.IsAbs(lpLibFileName) {
		return fmt.Errorf("lpLibFileName must be absolute")
	}

	dll, err := windows.LoadDLL(lpLibFileName)
	if err != nil {
		return fmt.Errorf("Loading DLL failed: %w", err)
	}

	defer func() {
		canUnloadProc, err := dll.FindProc("DllCanUnloadNow")
		if err != nil {
			return
		}

		if r1, _, _ := canUnloadProc.Call(); r1 != windows.NO_ERROR {
			return
		}

		dll.Release()
	}()

	createProc, err := dll.FindProc("CreateWebViewEnvironmentWithOptionsInternal")
	if err != nil {
		return fmt.Errorf("Unable to find CreateWebViewEnvironmentWithOptionsInternal entrypoint: %w", err)
	}

	userDataPtr, err := windows.UTF16PtrFromString(userDataFolder)
	if err != nil {
		return err
	}

	envOptionsCom := combridge.New2[iCoreWebView2EnvironmentOptions, iCoreWebView2EnvironmentOptions2](
		envOptions, envOptions)

	defer envOptionsCom.Close()

	envCompletedHandler = &environmentCreatedHandler{envCompletedHandler}
	envCompletedCom := combridge.New[iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler](envCompletedHandler)
	defer envCompletedCom.Close()

	preventEnvAndRegistryOverrides()

	const unknown = 1
	hr, _, err := createProc.Call(
		uintptr(unknown),
		uintptr(runtimeType),
		uintptr(unsafe.Pointer(userDataPtr)),
		uintptr(envOptionsCom.Ref()),
		uintptr(envCompletedCom.Ref()))

	if hr != 0 {
		if err == nil || err == windows.ERROR_SUCCESS {
			err = syscall.Errno(hr)
		}
		return err
	}

	return nil
}

type environmentCreatedHandler struct {
	originalHandler ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
}

func (r *environmentCreatedHandler) EnvironmentCompleted(errorCode HRESULT, createdEnvironment *ICoreWebView2Environment) HRESULT {
	// The OpenWebview2Loader has some retry logic and retries once, didn't encounter any case when this would have been
	// needed during the development: https://github.com/jchv/OpenWebView2Loader/blob/master/Source/WebView2Loader.cpp#L202

	if createdEnvironment != nil {
		// May or may not be necessary, but the official WebView2Loader seems to do it.
		iidICoreWebView2Environment := windows.GUID{
			Data1: 0xb96d755e,
			Data2: 0x0319,
			Data3: 0x4e92,
			Data4: [8]byte{0xa2, 0x96, 0x23, 0x43, 0x6f, 0x46, 0xa1, 0xfc},
		}

		if err := createdEnvironment.QueryInterface(&iidICoreWebView2Environment, &createdEnvironment); err != nil {
			createdEnvironment = nil
			errNo, ok := err.(syscall.Errno)
			if !ok {
				errNo = syscall.Errno(windows.E_FAIL)
			}
			errorCode = HRESULT(errNo)
		}
	}

	r.originalHandler.EnvironmentCompleted(errorCode, createdEnvironment)

	if createdEnvironment != nil {
		createdEnvironment.Release()
	}

	return HRESULT(windows.S_OK)
}

func preventEnvAndRegistryOverrides() {
	// Setting these env variables to empty string also prevents registry overrides because webview2
	// checks for existence and not for empty value
	os.Setenv("WEBVIEW2_PIPE_FOR_SCRIPT_DEBUGGER", "")
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "")
	os.Setenv("WEBVIEW2_RELEASE_CHANNEL_PREFERENCE", "0")

	// The following seems not be be required because those are only used by the webview2loader which
	// in this case is implemented on our own. But nevertheless set them to empty to be consistent.
	os.Setenv("WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", "")
	os.Setenv("WEBVIEW2_USER_DATA_FOLDER", "")
}
