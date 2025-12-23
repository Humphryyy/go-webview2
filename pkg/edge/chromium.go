//go:build windows
// +build windows

package edge

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/Humphryyy/go-webview2/internal/w32"
	"github.com/Humphryyy/go-webview2/webviewloader"
	"golang.org/x/sys/windows"
)

type Rect = w32.Rect

func globalErrorHandler(err error) {
	if err == nil {
		return
	}

	fmt.Printf("[WebView2 Error] %v\n", err)

	stackBuf := make([]uintptr, 64)
	stackSize := runtime.Callers(2, stackBuf)
	frames := runtime.CallersFrames(stackBuf[:stackSize])

	fmt.Println("\nStack trace:")
	stackIndex := 1
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		log.Printf("%d: %s\n\t%s:%d\n", stackIndex, frame.Function, frame.File, frame.Line)
		stackIndex++
	}
}

type Chromium struct {
	hwnd    uintptr
	padding struct {
		Left   int32
		Top    int32
		Right  int32
		Bottom int32
	}

	controller                       *ICoreWebView2Controller
	webview                          *ICoreWebView2
	inited                           uintptr
	envCompleted                     *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controllerCompleted              *iCoreWebView2CreateCoreWebView2ControllerCompletedHandler
	webMessageReceived               *iCoreWebView2WebMessageReceivedEventHandler
	containsFullScreenElementChanged *ICoreWebView2ContainsFullScreenElementChangedEventHandler
	permissionRequested              *iCoreWebView2PermissionRequestedEventHandler
	webResourceRequested             *iCoreWebView2WebResourceRequestedEventHandler
	acceleratorKeyPressed            *ICoreWebView2AcceleratorKeyPressedEventHandler
	navigationCompleted              *ICoreWebView2NavigationCompletedEventHandler
	processFailed                    *ICoreWebView2ProcessFailedEventHandler

	environment            *ICoreWebView2Environment
	webview2RuntimeVersion string

	// Settings
	Debug                 bool
	DataPath              string
	BrowserPath           string
	AdditionalBrowserArgs []string

	// permissions
	permissions      map[CoreWebView2PermissionKind]CoreWebView2PermissionState
	globalPermission *CoreWebView2PermissionState

	// Callbacks
	MessageCallback                          func(message string, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs)
	MessageWithAdditionalObjectsCallback     func(message string, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs)
	WebResourceRequestedCallback             func(request *ICoreWebView2WebResourceRequest, args *ICoreWebView2WebResourceRequestedEventArgs)
	WebResourceResponseReceivedCallback      func(response *ICoreWebView2WebResourceResponseView, args *ICoreWebView2WebResourceResponseReceivedEventArgs)
	NavigationCompletedCallback              func(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs)
	ProcessFailedCallback                    func(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs)
	ContainsFullScreenElementChangedCallback func(sender *ICoreWebView2, args *ICoreWebView2ContainsFullScreenElementChangedEventArgs)
	AcceleratorKeyCallback                   func(uint) bool

	// Error handling
	globalErrorCallback func(error)

	shuttingDown bool

	// Resize debouncing
	lastBounds  *w32.Rect
	resizeTimer *time.Timer
}

func NewChromium() *Chromium {
	e := &Chromium{}
	/*
	 All these handlers are passed to native code through syscalls with 'uintptr(unsafe.Pointer(handler))' and we know
	 that a pointer to those will be kept in the native code. Furthermore these handlers als contain pointer to other Go
	 structs like the vtable.
	 This violates the unsafe.Pointer rule '(4) Conversion of a Pointer to a uintptr when calling syscall.Syscall.' because
	 theres no guarantee that Go doesn't move these objects.
	 AFAIK currently the Go runtime doesn't move HEAP objects, so we should be safe with these handlers. But they don't
	 guarantee it, because in the future Go might use a compacting GC.
	 There's a proposal to add a runtime.Pin function, to prevent moving pinned objects, which would allow to easily fix
	 this issue by just pinning the handlers. The https://go-review.googlesource.com/c/go/+/367296/ should land in Go 1.19.
	*/
	e.envCompleted = newICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(e)
	e.controllerCompleted = newICoreWebView2CreateCoreWebView2ControllerCompletedHandler(e)
	e.webMessageReceived = newICoreWebView2WebMessageReceivedEventHandler(e)
	e.permissionRequested = newICoreWebView2PermissionRequestedEventHandler(e)
	e.webResourceRequested = newICoreWebView2WebResourceRequestedEventHandler(e)
	e.acceleratorKeyPressed = newICoreWebView2AcceleratorKeyPressedEventHandler(e)
	e.navigationCompleted = newICoreWebView2NavigationCompletedEventHandler(e)
	e.processFailed = newICoreWebView2ProcessFailedEventHandler(e)
	e.containsFullScreenElementChanged = newICoreWebView2ContainsFullScreenElementChangedEventHandler(e)
	/*
		// Pinner seems to panic in some cases as reported on Discord, maybe during shutdown when GC detects pinned objects
		// to be released that have not been unpinned.
		// It would also be better to use our ComBridge for this event handlers implementation instead of pinning them.
		// So all COM Implementations on the go-side use the same code.
		var pinner runtime.Pinner
		pinner.Pin(e.envCompleted)
		pinner.Pin(e.controllerCompleted)
		pinner.Pin(e.webMessageReceived)
		pinner.Pin(e.permissionRequested)
		pinner.Pin(e.webResourceRequested)
		pinner.Pin(e.acceleratorKeyPressed)
		pinner.Pin(e.navigationCompleted)
		pinner.Pin(e.processFailed)
		pinner.Pin(e.containsFullScreenElementChanged)
	*/
	e.permissions = make(map[CoreWebView2PermissionKind]CoreWebView2PermissionState)
	e.globalErrorCallback = globalErrorHandler
	return e
}

func (e *Chromium) ShuttingDown() {
	e.shuttingDown = true
}

func (e *Chromium) errorCallback(err error) {
	e.globalErrorCallback(err)
}

func (e *Chromium) SetErrorCallback(callback func(error)) {
	if callback != nil {
		e.globalErrorCallback = callback
	}
}

func (e *Chromium) Embed(hwnd uintptr) bool {

	var err error

	e.hwnd = hwnd

	dataPath := e.DataPath
	if dataPath == "" {
		currentExePath := make([]uint16, windows.MAX_PATH)
		_, err = windows.GetModuleFileName(windows.Handle(0), &currentExePath[0], windows.MAX_PATH)
		if err != nil {
			e.errorCallback(err)
		}
		currentExeName := filepath.Base(windows.UTF16ToString(currentExePath))
		dataPath = filepath.Join(os.Getenv("AppData"), currentExeName)
	}

	if e.BrowserPath != "" {
		if _, err = os.Stat(e.BrowserPath); errors.Is(err, os.ErrNotExist) {
			e.errorCallback(fmt.Errorf("browser path '%s' does not exist", e.BrowserPath))
		}
	}

	browserArgs := strings.Join(e.AdditionalBrowserArgs, " ")
	if err := createCoreWebView2EnvironmentWithOptions(e.BrowserPath, dataPath, e.envCompleted, browserArgs); err != nil {
		e.errorCallback(fmt.Errorf("error calling Webview2Loader: %s", err.Error()))
	}

	e.webview2RuntimeVersion, err = webviewloader.GetAvailableCoreWebView2BrowserVersionString(e.BrowserPath)
	if err != nil {
		e.errorCallback(fmt.Errorf("error getting Webview2 runtime version: %s", err.Error()))
	}

	var msg w32.Msg
	for {
		if atomic.LoadUintptr(&e.inited) != 0 {
			break
		}
		r, _, _ := w32.User32GetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if r == 0 {
			break
		}
		w32.User32TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		w32.User32DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	e.Init("window.external={invoke:s=>window.chrome.webview.postMessage(s)}")
	return true
}

func (e *Chromium) SetPadding(padding Rect) {
	if e.padding.Left == padding.Left && e.padding.Top == padding.Top &&
		e.padding.Right == padding.Right && e.padding.Bottom == padding.Bottom {

		return
	}

	e.padding.Left = padding.Left
	e.padding.Top = padding.Top
	e.padding.Right = padding.Right
	e.padding.Bottom = padding.Bottom
	e.Resize()
}

func (e *Chromium) ResizeWithBounds(bounds *Rect) {
	if e.hwnd == 0 {
		return
	}

	bounds.Top += e.padding.Top
	bounds.Bottom -= e.padding.Bottom
	bounds.Left += e.padding.Left
	bounds.Right -= e.padding.Right

	e.SetSize(*bounds)
}

func (e *Chromium) Resize() {
	if e.hwnd == 0 {
		return
	}

	bounds, err := w32.GetClientRect(e.hwnd)
	if err != nil {
		e.errorCallback(err)
		return
	}

	e.ResizeWithBounds(&bounds)
}

func (e *Chromium) Navigate(url string) {
	errStates := []string{"The group or resource is not in the correct state to perform the requested operation.", "This method can only be called from the thread that created the object."}

	maxRetries := 5
	delay := 100 * time.Millisecond

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := e.webview.Navigate(url)
		if err == nil {
			return
		}

		isKnown := false
		for _, errState := range errStates {
			if err.Error() == errState {
				isKnown = true
				break
			}
		}

		if !isKnown || attempt == maxRetries {
			e.errorCallback(err)
			return
		}

		time.Sleep(delay)
		delay *= 2
	}
}

func (e *Chromium) NavigateToString(content string) {
	err := e.webview.NavigateToString(content)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) Init(script string) {
	err := e.webview.AddScriptToExecuteOnDocumentCreated(script, nil)
	if err != nil {
		e.errorCallback(err)
	}
}

var (
	cdpHandlerPool     = make([]*ICoreWebView2CallDevToolsProtocolMethodCompletedHandler, 0)
	scriptHandlerPool  = make([]*ICoreWebView2ExecuteScriptCompletedHandler, 0)
	scriptResultPool   = make([]*ICoreWebView2ExecuteScriptWithResultCompletedHandler, 0)
	handlerPoolMutex   sync.Mutex
	maxHandlerPoolSize = 100000
)

type CallDevToolsProtocolMethodCompletedHandler struct {
	resultFunc func(errorCode uintptr, result string) uintptr
	refCount   int32
	poolKey    uintptr
	method     string
}

func (c *CallDevToolsProtocolMethodCompletedHandler) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (c *CallDevToolsProtocolMethodCompletedHandler) AddRef() uintptr {
	return uintptr(atomic.AddInt32(&c.refCount, 1))
}

func (c *CallDevToolsProtocolMethodCompletedHandler) Release() uintptr {
	newCount := atomic.AddInt32(&c.refCount, -1)
	if newCount == 0 && c.poolKey != 0 {
		// Remove from pool when WebView2 releases it
		handlerPoolMutex.Lock()
		for i, h := range cdpHandlerPool {
			if uintptr(unsafe.Pointer(h)) == c.poolKey {
				cdpHandlerPool = append(cdpHandlerPool[:i], cdpHandlerPool[i+1:]...)
				break
			}
		}
		handlerPoolMutex.Unlock()
	}
	return uintptr(newCount)
}

func (c *CallDevToolsProtocolMethodCompletedHandler) CallDevToolsProtocolMethodCompleted(errorCode uintptr, result string) uintptr {
	return c.resultFunc(errorCode, result)
}

func (e *Chromium) CallDevToolsProtocolMethod(method string, params string, callback func(errorCode uintptr, result string) uintptr) error {
	if e.webview == nil {
		return errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return errors.New("WebView2 runtime version not available")
	}

	handlerImpl := &CallDevToolsProtocolMethodCompletedHandler{
		resultFunc: callback,
		refCount:   1, // Start with refcount of 1
		method:     method,
	}

	handler := NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(handlerImpl)
	handlerImpl.poolKey = uintptr(unsafe.Pointer(handler))

	// Keep a reference to prevent garbage collection
	handlerPoolMutex.Lock()
	cdpHandlerPool = append(cdpHandlerPool, handler)
	if len(cdpHandlerPool) > maxHandlerPoolSize {
		cdpHandlerPool = cdpHandlerPool[len(cdpHandlerPool)/2:]
	}
	handlerPoolMutex.Unlock()

	err := e.webview.CallDevToolsProtocolMethod(method, params, handler)
	if err != nil {
		return fmt.Errorf("error calling dev tools protocol method: %w", err)
	}

	return nil
}

func (e *Chromium) CallDevToolsProtocolMethodForSession(sessionId string, method string, params string, callback func(errorCode uintptr, result string) uintptr) error {
	if e.webview == nil {
		return errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return errors.New("WebView2 runtime version not available")
	}

	handlerImpl := &CallDevToolsProtocolMethodCompletedHandler{
		resultFunc: callback,
		refCount:   1, // Start with refcount of 1
		method:     method,
	}

	handler := NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(handlerImpl)
	handlerImpl.poolKey = uintptr(unsafe.Pointer(handler))

	// Keep a reference to prevent garbage collection
	handlerPoolMutex.Lock()
	cdpHandlerPool = append(cdpHandlerPool, handler)
	if len(cdpHandlerPool) > maxHandlerPoolSize {
		cdpHandlerPool = cdpHandlerPool[len(cdpHandlerPool)/2:]
	}
	handlerPoolMutex.Unlock()

	// Get ICoreWebView2_11 interface
	webview2 := e.webview.GetICoreWebView2_11()
	if webview2 == nil {
		return errors.New("failed to get ICoreWebView2_11 interface from webview")
	}

	defer webview2.Release()

	err := webview2.CallDevToolsProtocolMethodForSession(sessionId, method, params, handler)
	if err != nil {
		return fmt.Errorf("error calling dev tools protocol method for session: %w", err)
	}

	return nil
}

func (e *Chromium) GetDevToolsProtocolEventReceiver(eventName string) (*ICoreWebView2DevToolsProtocolEventReceiver, error) {
	if e.webview == nil {
		return nil, errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return nil, errors.New("WebView2 runtime version not available")
	}

	return e.webview.GetDevToolsProtocolEventReceiver(eventName)
}

type ExecuteScriptCompletedHandler struct {
	resultFunc func(errorCode uintptr, executedScript string) uintptr
	refCount   int32
	poolKey    uintptr
}

func (e *ExecuteScriptCompletedHandler) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (e *ExecuteScriptCompletedHandler) AddRef() uintptr {
	return uintptr(atomic.AddInt32(&e.refCount, 1))
}

func (e *ExecuteScriptCompletedHandler) Release() uintptr {
	newCount := atomic.AddInt32(&e.refCount, -1)
	if newCount == 0 && e.poolKey != 0 {
		// Remove from pool when WebView2 releases it
		handlerPoolMutex.Lock()
		for i, h := range scriptHandlerPool {
			if uintptr(unsafe.Pointer(h)) == e.poolKey {
				scriptHandlerPool = append(scriptHandlerPool[:i], scriptHandlerPool[i+1:]...)
				break
			}
		}
		handlerPoolMutex.Unlock()
	}
	return uintptr(newCount)
}

func (e *ExecuteScriptCompletedHandler) ExecuteScriptCompleted(errorCode uintptr, executedScript string) uintptr {
	return e.resultFunc(errorCode, executedScript)
}

func (e *Chromium) Eval(script string, resultFunc ...func(errorCode uintptr, executedScript string) uintptr) error {
	if e.webview == nil || e.shuttingDown {
		return errors.New("webview not initialized or shutting down")
	}

	var handler *ICoreWebView2ExecuteScriptCompletedHandler

	if len(resultFunc) > 0 {
		handlerImpl := &ExecuteScriptCompletedHandler{
			resultFunc: resultFunc[0],
			refCount:   1, // Start with refcount of 1
		}

		handler = NewICoreWebView2ExecuteScriptCompletedHandler(handlerImpl)
		handlerImpl.poolKey = uintptr(unsafe.Pointer(handler))

		// Keep a reference to prevent garbage collection
		handlerPoolMutex.Lock()
		scriptHandlerPool = append(scriptHandlerPool, handler)
		if len(scriptHandlerPool) > maxHandlerPoolSize {
			scriptHandlerPool = scriptHandlerPool[len(scriptHandlerPool)/2:]
		}
		handlerPoolMutex.Unlock()
	}

	err := e.webview.ExecuteScript(script, handler)
	if err != nil {
		return fmt.Errorf("error executing script: %w", err)
	}

	return nil
}

type ExecuteScriptWithResultCompletedHandler struct {
	resultFunc func(errorCode uintptr, result *ICoreWebView2ExecuteScriptResult) uintptr
	refCount   int32
	poolKey    uintptr
}

func (e *ExecuteScriptWithResultCompletedHandler) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (e *ExecuteScriptWithResultCompletedHandler) AddRef() uintptr {
	return uintptr(atomic.AddInt32(&e.refCount, 1))
}

func (e *ExecuteScriptWithResultCompletedHandler) Release() uintptr {
	newCount := atomic.AddInt32(&e.refCount, -1)
	if newCount == 0 && e.poolKey != 0 {
		// Remove from pool when WebView2 releases it
		handlerPoolMutex.Lock()
		for i, h := range scriptResultPool {
			if uintptr(unsafe.Pointer(h)) == e.poolKey {
				scriptResultPool = append(scriptResultPool[:i], scriptResultPool[i+1:]...)
				break
			}
		}
		handlerPoolMutex.Unlock()
	}
	return uintptr(newCount)
}

func (e *ExecuteScriptWithResultCompletedHandler) ExecuteScriptWithResultCompleted(errorCode uintptr, result *ICoreWebView2ExecuteScriptResult) uintptr {
	return e.resultFunc(errorCode, result)
}

func (e *Chromium) EvalWithResult(script string, resultFunc func(errorCode uintptr, result *ICoreWebView2ExecuteScriptResult) uintptr) error {
	if e.webview == nil {
		return errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return errors.New("WebView2 runtime version not available")
	}

	handlerImpl := &ExecuteScriptWithResultCompletedHandler{
		resultFunc: resultFunc,
		refCount:   1,
	}

	handler := NewICoreWebView2ExecuteScriptWithResultCompletedHandler(handlerImpl)
	handlerImpl.poolKey = uintptr(unsafe.Pointer(handler))

	// Keep a reference to prevent garbage collection
	handlerPoolMutex.Lock()
	scriptResultPool = append(scriptResultPool, handler)
	if len(scriptResultPool) > maxHandlerPoolSize {
		scriptResultPool = scriptResultPool[len(scriptResultPool)/2:]
	}
	handlerPoolMutex.Unlock()

	// Get ICoreWebView2_21 interface
	webview2 := e.webview.GetICoreWebView2_21()
	if webview2 == nil {
		return errors.New("failed to get ICoreWebView2_21 interface from webview")
	}

	defer webview2.Release()

	err := webview2.ExecuteScriptWithResult(script, handler)
	if err != nil {
		return fmt.Errorf("error executing script with result: %w", err)
	}

	return nil
}

func (e *Chromium) Show() error {
	return e.controller.PutIsVisible(true)
}

func (e *Chromium) Hide() error {
	return e.controller.PutIsVisible(false)
}

func (e *Chromium) QueryInterface(_, _ uintptr) uintptr {
	return 0
}

func (e *Chromium) AddRef() uintptr {
	return 1
}

func (e *Chromium) Release() uintptr {
	return 1
}

func (e *Chromium) EnvironmentCompleted(res uintptr, env *ICoreWebView2Environment) uintptr {
	if env == nil {
		err := syscall.Errno(res)
		log.Printf("[WebView2] Environment creation failed with error code %v: %v\n", res, err)
		if e.globalErrorCallback != nil {
			e.globalErrorCallback(fmt.Errorf("failed to create WebView2 environment: %w", err))
		}
		return res
	}

	log.Printf("[WebView2] Environment created successfully\n")

	env.vtbl.AddRef.Call(uintptr(unsafe.Pointer(env)))
	e.environment = env

	err := env.CreateCoreWebView2Controller(e.hwnd, e.controllerCompleted)
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) CreateCoreWebView2ControllerCompleted(res uintptr, controller *ICoreWebView2Controller) uintptr {
	if int32(res) < 0 {
		e.errorCallback(fmt.Errorf("error creating controller with %08x: %s", res, syscall.Errno(res)))
	}

	var err error

	controller.vtbl.AddRef.Call(uintptr(unsafe.Pointer(controller)))
	e.controller = controller

	// Try to get ICoreWebView2Controller3 interface for better performance
	if controller3 := e.controller.GetICoreWebView2Controller3(); controller3 != nil {
		// Use raw pixels mode for better performance during resize
		if err := controller3.PutBoundsMode(COREWEBVIEW2_BOUNDS_MODE_USE_RAW_PIXELS); err != nil {
			e.errorCallback(err)
		}

		// Disable monitor scale changes since we're using raw pixels
		if err := controller3.PutShouldDetectMonitorScaleChanges(false); err != nil {
			e.errorCallback(err)
		}
	}
	var token _EventRegistrationToken
	e.webview, err = e.controller.GetCoreWebView2()
	if err != nil {
		e.errorCallback(err)
	}

	e.webview.vtbl.AddRef.Call(uintptr(unsafe.Pointer(e.webview)))
	err = e.webview.AddWebMessageReceived(e.webMessageReceived, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddPermissionRequested(e.permissionRequested, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddWebResourceRequested(e.webResourceRequested, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddNavigationCompleted(e.navigationCompleted, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddProcessFailed(e.processFailed, &token)
	if err != nil {
		e.errorCallback(err)
	}
	err = e.webview.AddContainsFullScreenElementChanged(e.containsFullScreenElementChanged, &token)
	if err != nil {
		e.errorCallback(err)
	}

	err = e.controller.AddAcceleratorKeyPressed(e.acceleratorKeyPressed, &token)
	if err != nil {
		e.errorCallback(err)
	}

	atomic.StoreUintptr(&e.inited, 1)

	return 0
}

func (e *Chromium) ContainsFullScreenElementChanged(sender *ICoreWebView2, args *ICoreWebView2ContainsFullScreenElementChangedEventArgs) uintptr {
	if e.ContainsFullScreenElementChangedCallback != nil {
		e.ContainsFullScreenElementChangedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) MessageReceived(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	message, err := args.TryGetWebMessageAsString()
	if err != nil {
		e.errorCallback(err)
	}

	if HasCapability(e.webview2RuntimeVersion, GetAdditionalObjects) {
		obj, err := args.GetAdditionalObjects()
		if err != nil {
			e.errorCallback(err)
		}

		if obj != nil && e.MessageWithAdditionalObjectsCallback != nil {
			defer obj.Release()
			e.MessageWithAdditionalObjectsCallback(message, sender, args)
		} else if e.MessageCallback != nil {
			e.MessageCallback(message, sender, args)
		}
	} else if e.MessageCallback != nil {
		e.MessageCallback(message, sender, args)
	}

	err = sender.PostWebMessageAsString(message)
	if err != nil && !errors.Is(err, windows.ERROR_IO_PENDING) {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) SetPermission(kind CoreWebView2PermissionKind, state CoreWebView2PermissionState) {
	e.permissions[kind] = state
}

func (e *Chromium) SetBackgroundColour(R, G, B, A uint8) {
	controller := e.GetController()
	controller2 := controller.GetICoreWebView2Controller2()

	backgroundCol := COREWEBVIEW2_COLOR{
		A: A,
		R: R,
		G: G,
		B: B,
	}

	// WebView2 only has 0 and 255 as valid values.
	if backgroundCol.A > 0 && backgroundCol.A < 255 {
		backgroundCol.A = 255
	}

	err := controller2.PutDefaultBackgroundColor(backgroundCol)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) SetGlobalPermission(state CoreWebView2PermissionState) {
	e.globalPermission = &state
}

func (e *Chromium) PermissionRequested(_ *ICoreWebView2, args *iCoreWebView2PermissionRequestedEventArgs) uintptr {
	kind, err := args.GetPermissionKind()
	if err != nil {
		e.errorCallback(err)
	}
	var result CoreWebView2PermissionState
	if e.globalPermission != nil {
		result = *e.globalPermission
	} else {
		var ok bool
		result, ok = e.permissions[kind]
		if !ok {
			result = CoreWebView2PermissionStateDefault
		}
	}
	err = args.PutState(result)
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) WebResourceRequested(sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr {
	req, err := args.GetRequest()
	if err != nil {
		log.Fatal(err)
	}
	defer req.Release()

	if e.WebResourceRequestedCallback != nil {
		e.WebResourceRequestedCallback(req, args)
	}
	return 0
}

func (e *Chromium) WebResourceResponseReceived(sender *ICoreWebView2, args *ICoreWebView2WebResourceResponseReceivedEventArgs) uintptr {
	resp, err := args.GetResponse()
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Release()

	if e.WebResourceResponseReceivedCallback != nil {
		e.WebResourceResponseReceivedCallback(resp, args)
	}

	return 0
}

func (e *Chromium) AddWebResourceRequestedFilter(filter string, ctx COREWEBVIEW2_WEB_RESOURCE_CONTEXT) {
	err := e.webview.AddWebResourceRequestedFilter(filter, ctx)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) Environment() *ICoreWebView2Environment {
	return e.environment
}

// AcceleratorKeyPressed is called when an accelerator key is pressed.
// If the AcceleratorKeyCallback method has been set, it will defer handling of the keypress
// to the callback. That callback returns a bool indicating if the event was handled.
func (e *Chromium) AcceleratorKeyPressed(sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr {
	if e.AcceleratorKeyCallback == nil {
		return 0
	}
	eventKind, _ := args.GetKeyEventKind()
	if eventKind == COREWEBVIEW2_KEY_EVENT_KIND_KEY_DOWN ||
		eventKind == COREWEBVIEW2_KEY_EVENT_KIND_SYSTEM_KEY_DOWN {
		virtualKey, _ := args.GetVirtualKey()
		status, _ := args.GetPhysicalKeyStatus()
		if !status.WasKeyDown {
			err := args.PutHandled(e.AcceleratorKeyCallback(virtualKey))
			if err != nil {
				e.errorCallback(err)
			}
		} else {
			return 0
		}
	}
	err := args.PutHandled(false)
	if err != nil {
		e.errorCallback(err)
	}
	return 0
}

func (e *Chromium) GetSettings() (*ICoreWebViewSettings, error) {
	return e.webview.GetSettings()
}

func (e *Chromium) GetController() *ICoreWebView2Controller {
	return e.controller
}

func boolToInt(input bool) int {
	if input {
		return 1
	}
	return 0
}

func (e *Chromium) NavigationCompleted(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	if e.NavigationCompletedCallback != nil {
		e.NavigationCompletedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) ProcessFailed(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr {
	if e.ProcessFailedCallback != nil {
		e.ProcessFailedCallback(sender, args)
	}
	return 0
}

func (e *Chromium) NotifyParentWindowPositionChanged() error {
	//It looks like the wndproc function is called before the controller initialization is complete.
	//Because of this the controller is nil
	if e.controller == nil {
		return nil
	}
	return e.controller.NotifyParentWindowPositionChanged()
}

func (e *Chromium) Focus() {
	err := e.controller.MoveFocus(COREWEBVIEW2_MOVE_FOCUS_REASON_PROGRAMMATIC)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) PutZoomFactor(zoomFactor float64) {
	err := e.controller.PutZoomFactor(zoomFactor)
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) OpenDevToolsWindow() {
	err := e.webview.OpenDevToolsWindow()
	if err != nil {
		e.errorCallback(err)
	}
}

func (e *Chromium) HasCapability(c Capability) bool {
	return HasCapability(e.webview2RuntimeVersion, c)
}

func (e *Chromium) GetIsSwipeNavigationEnabled() (bool, error) {
	if !HasCapability(e.webview2RuntimeVersion, SwipeNavigation) {
		return false, UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return false, err
	}
	webview2Settings6 := webview2Settings.GetICoreWebView2Settings6()
	var result bool
	result, err = webview2Settings6.GetIsSwipeNavigationEnabled()
	if err != nil {
		return false, err
	}
	return result, nil
}

// PutIsGeneralAutofillEnabled controls whether autofill for information
// like names, street and email addresses, phone numbers, and arbitrary input
// is enabled. This excludes password and credit card information. When
// IsGeneralAutofillEnabled is false, no suggestions appear, and no new information
// is saved. When IsGeneralAutofillEnabled is true, information is saved, suggestions
// appear and clicking on one will populate the form fields.
// It will take effect immediately after setting.
// The default value is `FALSE`.
func (e *Chromium) PutIsGeneralAutofillEnabled(value bool) error {
	if !HasCapability(e.webview2RuntimeVersion, GeneralAutofillEnabled) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings4 := webview2Settings.GetICoreWebView2Settings4()
	return webview2Settings4.PutIsGeneralAutofillEnabled(value)
}

// PutIsPasswordAutosaveEnabled sets whether the browser should offer to save passwords and other
// identifying information entered into forms automatically.
// The default value is `FALSE`.
func (e *Chromium) PutIsPasswordAutosaveEnabled(value bool) error {
	if !HasCapability(e.webview2RuntimeVersion, PasswordAutosaveEnabled) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings4 := webview2Settings.GetICoreWebView2Settings4()
	return webview2Settings4.PutIsPasswordAutosaveEnabled(value)
}

func (e *Chromium) PutIsSwipeNavigationEnabled(enabled bool) error {
	if !HasCapability(e.webview2RuntimeVersion, SwipeNavigation) {
		return UnsupportedCapabilityError
	}
	webview2Settings, err := e.webview.GetSettings()
	if err != nil {
		return err
	}
	webview2Settings6 := webview2Settings.GetICoreWebView2Settings6()
	err = webview2Settings6.PutIsSwipeNavigationEnabled(enabled)
	if err != nil {
		return err
	}
	return nil
}

func (e *Chromium) AllowExternalDrag(allow bool) error {
	if !HasCapability(e.webview2RuntimeVersion, AllowExternalDrop) {
		return UnsupportedCapabilityError
	}
	controller := e.GetController()
	controller4 := controller.GetICoreWebView2Controller4()
	err := controller4.PutAllowExternalDrop(allow)
	if err != nil {
		return err
	}
	return nil
}

func (e *Chromium) GetAllowExternalDrag() (bool, error) {
	if !HasCapability(e.webview2RuntimeVersion, AllowExternalDrop) {
		return false, UnsupportedCapabilityError
	}
	controller := e.GetController()
	controller4 := controller.GetICoreWebView2Controller4()
	result, err := controller4.GetAllowExternalDrop()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (e *Chromium) GetCookieManager() (*ICoreWebView2CookieManager, error) {
	if e.webview == nil {
		return nil, errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return nil, errors.New("WebView2 runtime version not available")
	}

	// Get ICoreWebView2_2 interface
	webview2, err := e.webview.QueryInterface2()
	if err != nil {
		return nil, fmt.Errorf("failed to get ICoreWebView2_2: %w\nThis functionality requires WebView2 Runtime version 89.0.721.0 or later. Current version: %s", err, e.webview2RuntimeVersion)
	}
	defer webview2.Release()

	// Get cookie manager
	cookieManager, err := webview2.GetCookieManager()
	if err != nil {
		return nil, fmt.Errorf("failed to get cookie manager: %w", err)
	}

	// Note: The caller is responsible for calling Release() on the returned cookieManager
	return cookieManager, nil
}

func (e *Chromium) AddWebResourceRequestedFilterWithRequestSourceKinds(uri string, ResourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT, requestSourceKinds COREWEBVIEW2_WEB_RESOURCE_REQUEST_SOURCE_KINDS) error {
	if e.webview == nil {
		return errors.New("webview not initialized")
	}

	// Check WebView2 version
	if e.webview2RuntimeVersion == "" {
		return errors.New("WebView2 runtime version not available")
	}

	// Get ICoreWebView2_22 interface
	webview2 := e.webview.GetICoreWebView2_22()
	if webview2 == nil {
		return errors.New("failed to get ICoreWebView2_22 interface from webview")
	}

	defer webview2.Release()

	// Add web resource requested filter
	err := webview2.AddWebResourceRequestedFilterWithRequestSourceKinds(uri, ResourceContext, requestSourceKinds)
	if err != nil {
		return fmt.Errorf("failed to add web resource requested filter: %w", err)
	}

	return nil
}
