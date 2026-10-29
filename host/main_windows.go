//go:build windows

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	wmUser         = 0x0400
	wmClose        = 0x0010
	wmCloseReady   = wmUser + 8
	wmCommand      = 0x0111
	wmDestroy      = 0x0002
	imageIcon      = 1
	lrLoadFromFile = 0x00000010
	idOpen         = 1
	idQuit         = 2
	mbIconError    = 0x00000010
	swShowNormal   = 1
	swHide         = 0
	wmSize         = 0x0005
	wsOverlapped   = 0x00CF0000
	className      = "TurnipProphetFrame"
)

var (
	user32  = windows.NewLazySystemDLL("user32.dll")
	shell32 = windows.NewLazySystemDLL("shell32.dll")

	procRegisterClassExW         = user32.NewProc("RegisterClassExW")
	procCreateWindowExW          = user32.NewProc("CreateWindowExW")
	procDefWindowProcW           = user32.NewProc("DefWindowProcW")
	procGetMessageW              = user32.NewProc("GetMessageW")
	procTranslateMessage         = user32.NewProc("TranslateMessage")
	procDispatchMessageW         = user32.NewProc("DispatchMessageW")
	procPostQuitMessage          = user32.NewProc("PostQuitMessage")
	procPostMessageW             = user32.NewProc("PostMessageW")
	procLoadImageW               = user32.NewProc("LoadImageW")
	procLoadIconW                = user32.NewProc("LoadIconW")
	procSetForegroundWindow      = user32.NewProc("SetForegroundWindow")
	procMessageBoxW              = user32.NewProc("MessageBoxW")
	procShellExecuteW            = shell32.NewProc("ShellExecuteW")
	procFindWindowW              = user32.NewProc("FindWindowW")
	procShowWindow               = user32.NewProc("ShowWindow")
	procDestroyWindow            = user32.NewProc("DestroyWindow")
	procGetSystemMetrics         = user32.NewProc("GetSystemMetrics")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procAllowSetForegroundWindow = user32.NewProc("AllowSetForegroundWindow")
	procSetAppID                 = shell32.NewProc("SetCurrentProcessExplicitAppUserModelID")

	frameWnd     windows.HWND
	appPort      int
	appWork      string
	appRoot      string
	closeStarted bool
	exitSaved    = make(chan struct{})
	exitOnce     sync.Once
)

type point struct {
	X int32
	Y int32
}

type msg struct {
	HWND    windows.HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      point
}

type wndClassEx struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   windows.Handle
	Icon       windows.Handle
	Cursor     windows.Handle
	Background windows.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     windows.Handle
}

var wndProcCB = syscall.NewCallback(wndProc)

func utf16Ptr(s string) *uint16 {
	p, err := windows.UTF16PtrFromString(s)
	if err != nil {
		return nil
	}
	return p
}

func showError(text string) {
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(utf16Ptr(text))), uintptr(unsafe.Pointer(utf16Ptr("Turnip Prophet"))), mbIconError)
}

func exeDir() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	p, err = filepath.EvalSymlinks(p)
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

func findBrowser() string {
	names := []string{
		filepath.Join("Microsoft", "Edge", "Application", "msedge.exe"),
		filepath.Join("Google", "Chrome", "Application", "chrome.exe"),
	}
	roots := []string{os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)"), os.Getenv("LocalAppData")}
	for _, root := range roots {
		if root == "" {
			continue
		}
		for _, name := range names {
			p := filepath.Join(root, name)
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p
			}
		}
	}
	return ""
}

func openBrowser() {
	if appPort == 0 {
		return
	}
	url := fmt.Sprintf("http://127.0.0.1:%d/", appPort)
	browser := findBrowser()
	if browser == "" {
		procShellExecuteW.Call(0, 0, uintptr(unsafe.Pointer(utf16Ptr(url))), 0, 0, swShowNormal)
		return
	}
	profile := filepath.Join(appWork, "profile")
	_ = os.MkdirAll(profile, 0700)
	cmd := exec.Command(browser,
		"--user-data-dir="+profile,
		"--app="+url,
		"--window-size=1360,960",
		"--no-first-run",
		"--disable-sync",
	)
	_ = cmd.Start()
}

func loadIcon() windows.Handle {
	path := filepath.Join(appRoot, "TurnipProphet.ico")
	h, _, _ := procLoadImageW.Call(0, uintptr(unsafe.Pointer(utf16Ptr(path))), imageIcon, 0, 0, lrLoadFromFile)
	if h != 0 {
		return windows.Handle(h)
	}
	h, _, _ = procLoadIconW.Call(0, 32512)
	return windows.Handle(h)
}

func signalExitSaved() {
	exitOnce.Do(func() { close(exitSaved) })
}

func beginClose(hwnd uintptr) {
	if closeStarted {
		return
	}
	closeStarted = true
	if browser != nil {
		browser.Eval(`window.flushSaveAndExit&&window.flushSaveAndExit()`)
	} else {
		signalExitSaved()
	}
	go func() {
		select {
		case <-exitSaved:
		case <-time.After(1500 * time.Millisecond):
		}
		procPostMessageW.Call(hwnd, wmCloseReady, 0, 0)
	}()
}

func wndProc(hwnd, message, wparam, lparam uintptr) uintptr {
	switch message {
	case wmClose:
		beginClose(hwnd)
		return 0
	case wmCloseReady:
		procDestroyWindow.Call(hwnd)
		return 0
	case wmSize:
		if browser != nil {
			browser.Resize()
		}
		r, _, _ := procDefWindowProcW.Call(hwnd, message, wparam, lparam)
		return r
	case wmCommand:
		switch uint16(wparam) {
		case idOpen:
			showFrame()
		case idQuit:
			procPostQuitMessage.Call(0)
		}
		return 0
	case wmDestroy:
		procPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := procDefWindowProcW.Call(hwnd, message, wparam, lparam)
	return r
}

func postQuit() {
	if frameWnd != 0 {
		procPostMessageW.Call(uintptr(frameWnd), wmCommand, idQuit, 0)
	}
}

func runLoop() {
	var m msg
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(ret) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func acquireMutex() (windows.Handle, bool) {
	name := utf16Ptr(`Local\TurnipProphet`)
	h, err := windows.CreateMutex(nil, false, name)
	return h, err == windows.ERROR_ALREADY_EXISTS
}

func setAppID() {
	procSetAppID.Call(uintptr(unsafe.Pointer(utf16Ptr("foctober.TurnipProphet"))))
}

func signalRunning() bool {
	hwnd, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(utf16Ptr(className))), 0)
	if hwnd == 0 {
		return false
	}
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	if pid != 0 {
		procAllowSetForegroundWindow.Call(uintptr(pid))
	}
	procPostMessageW.Call(hwnd, wmCommand, idOpen, 0)
	return true
}

func showFrame() {
	if frameWnd == 0 {
		openBrowser()
		return
	}
	procShowWindow.Call(uintptr(frameWnd), swShowNormal)
	procSetForegroundWindow.Call(uintptr(frameWnd))
	setTitleBarTheme(frameWnd)
	if browser != nil {
		_ = browser.Show()
		browser.Resize()
		browser.Focus()
	}
}

func main() {
	setAppID()
	root, err := exeDir()
	if err != nil {
		showError(err.Error())
		return
	}
	appRoot = root
	if miss := MissingFiles(root); len(miss) > 0 {
		showError("Missing " + miss[0] + ". Reinstall Turnip Prophet from TurnipProphet-Setup.exe.")
		return
	}
	if port := FindRunning(); port > 0 {
		for i := 0; i < 30 && !signalRunning(); i++ {
			time.Sleep(100 * time.Millisecond)
		}
		return
	}
	mutex, already := acquireMutex()
	if mutex != 0 {
		defer windows.CloseHandle(mutex)
	}
	if already {
		for i := 0; i < 30; i++ {
			if signalRunning() {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		showError("Turnip Prophet is already starting. Give it a second, then open the shortcut again.")
		return
	}

	app, err := Prepare(root)
	if err != nil {
		showError(err.Error())
		return
	}
	appWork = app.Work
	app.onQuit = postQuit
	app.onExitSaved = signalExitSaved

	ln, port, err := Listen()
	if err != nil {
		showError(err.Error() + " Quit old copies and try again.")
		return
	}
	appPort = port
	writePortFile(app.Work, port)
	defer clearPortFile(app.Work)
	app.logf("Listening on 127.0.0.1:%d", port)

	srv := &http.Server{Handler: app.Handler(), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			app.logf("server: %v", err)
		}
	}()
	defer srv.Shutdown(context.Background())

	hinst, _, _ := windows.NewLazySystemDLL("kernel32.dll").NewProc("GetModuleHandleW").Call(0)
	class := utf16Ptr(className)
	icon := loadIcon()
	wc := wndClassEx{
		WndProc:   wndProcCB,
		Instance:  windows.Handle(hinst),
		ClassName: class,
		Icon:      icon,
		IconSm:    icon,
	}
	wc.Size = uint32(unsafe.Sizeof(wc))
	if atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		showError("Could not create the window. " + err.Error())
		return
	}
	screenW, _, _ := procGetSystemMetrics.Call(0)
	screenH, _, _ := procGetSystemMetrics.Call(1)
	const winW, winH = 1360, 960
	posX := (int32(screenW) - winW) / 2
	posY := (int32(screenH) - winH) / 2
	if posX < 0 {
		posX = 0
	}
	if posY < 0 {
		posY = 0
	}
	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(class)),
		uintptr(unsafe.Pointer(utf16Ptr("Turnip Prophet"))),
		wsOverlapped,
		uintptr(posX), uintptr(posY), winW, winH,
		0, 0, hinst, 0,
	)
	if hwnd == 0 {
		showError("Could not create the window. " + err.Error())
		return
	}
	frameWnd = windows.HWND(hwnd)
	stampWindow(frameWnd)
	procShowWindow.Call(hwnd, swShowNormal)
	attachPage()
	runLoop()
	app.logf("Stopped.")
}
