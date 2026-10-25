//go:build windows

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	wmUser          = 0x0400
	wmTray          = wmUser + 1
	wmNull          = 0x0000
	wmCommand       = 0x0111
	wmDestroy       = 0x0002
	wmLButtonDblClk = 0x0203
	wmRButtonUp     = 0x0205
	nimAdd          = 0x00000000
	nimDelete       = 0x00000002
	nifMessage      = 0x00000001
	nifIcon         = 0x00000002
	nifTip          = 0x00000004
	nifInfo         = 0x00000010
	niifInfo        = 0x00000001
	imageIcon       = 1
	lrLoadFromFile  = 0x00000010
	mfString        = 0x00000000
	tpmRightAlign   = 0x0008
	tpmBottomAlign  = 0x0020
	tpmRightButton  = 0x0002
	idOpen          = 1
	idQuit          = 2
	mbIconError     = 0x00000010
	swShowNormal    = 1
	className       = "TurnipProphetHidden"
)

var (
	user32  = windows.NewLazySystemDLL("user32.dll")
	shell32 = windows.NewLazySystemDLL("shell32.dll")

	procRegisterClassExW    = user32.NewProc("RegisterClassExW")
	procCreateWindowExW     = user32.NewProc("CreateWindowExW")
	procDefWindowProcW      = user32.NewProc("DefWindowProcW")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")
	procPostQuitMessage     = user32.NewProc("PostQuitMessage")
	procPostMessageW        = user32.NewProc("PostMessageW")
	procLoadImageW          = user32.NewProc("LoadImageW")
	procLoadIconW           = user32.NewProc("LoadIconW")
	procCreatePopupMenu     = user32.NewProc("CreatePopupMenu")
	procAppendMenuW         = user32.NewProc("AppendMenuW")
	procTrackPopupMenu      = user32.NewProc("TrackPopupMenu")
	procDestroyMenu         = user32.NewProc("DestroyMenu")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procGetCursorPos        = user32.NewProc("GetCursorPos")
	procShellNotifyIconW    = shell32.NewProc("Shell_NotifyIconW")
	procMessageBoxW         = user32.NewProc("MessageBoxW")
	procShellExecuteW       = shell32.NewProc("ShellExecuteW")

	trayWnd windows.HWND
	trayID  uint32 = 1
	appPort int
	appWork string
	appRoot string
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

type notifyIconData struct {
	Size            uint32
	Wnd             windows.HWND
	ID              uint32
	Flags           uint32
	CallbackMessage uint32
	Icon            windows.Handle
	Tip             [128]uint16
	State           uint32
	StateMask       uint32
	Info            [256]uint16
	Version         uint32
	InfoTitle       [64]uint16
	InfoFlags       uint32
	GuidItem        [16]byte
	BalloonIcon     windows.Handle
}

const (
	_ uintptr = 976 - unsafe.Sizeof(notifyIconData{})
	_ uintptr = unsafe.Sizeof(notifyIconData{}) - 976
)

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
		"--window-size=1360,900",
		"--no-first-run",
		"--disable-sync",
	)
	_ = cmd.Start()
	go watchAppWindow()
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

func copyUTF16(dst []uint16, s string) {
	u := windows.StringToUTF16(s)
	n := len(u)
	if n > len(dst) {
		n = len(dst)
	}
	copy(dst[:n], u[:n])
}

func addTray(icon windows.Handle) {
	nid := notifyIconData{Size: uint32(unsafe.Sizeof(notifyIconData{})), Wnd: trayWnd, ID: trayID}
	nid.Flags = nifMessage | nifIcon | nifTip | nifInfo
	nid.CallbackMessage = wmTray
	nid.Icon = icon
	nid.InfoFlags = niifInfo
	copyUTF16(nid.Tip[:], "Turnip Prophet")
	copyUTF16(nid.InfoTitle[:], "Turnip Prophet")
	copyUTF16(nid.Info[:], "Tracker is open. Right-click this icon to quit.")
	procShellNotifyIconW.Call(nimAdd, uintptr(unsafe.Pointer(&nid)))
}

func removeTray() {
	nid := notifyIconData{Size: uint32(unsafe.Sizeof(notifyIconData{})), Wnd: trayWnd, ID: trayID}
	procShellNotifyIconW.Call(nimDelete, uintptr(unsafe.Pointer(&nid)))
}

func showMenu() {
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)
	procAppendMenuW.Call(menu, mfString, idOpen, uintptr(unsafe.Pointer(utf16Ptr("Open"))))
	procAppendMenuW.Call(menu, mfString, idQuit, uintptr(unsafe.Pointer(utf16Ptr("Quit"))))
	var pt point
	procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	procSetForegroundWindow.Call(uintptr(trayWnd))
	procTrackPopupMenu.Call(menu, tpmRightAlign|tpmBottomAlign|tpmRightButton, uintptr(pt.X), uintptr(pt.Y), 0, uintptr(trayWnd), 0)
	procPostMessageW.Call(uintptr(trayWnd), wmNull, 0, 0)
}

func wndProc(hwnd, message, wparam, lparam uintptr) uintptr {
	switch message {
	case wmTray:
		switch lparam {
		case wmLButtonDblClk:
			openBrowser()
		case wmRButtonUp:
			showMenu()
		}
		return 0
	case wmCommand:
		switch uint16(wparam) {
		case idOpen:
			openBrowser()
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
	if trayWnd != 0 {
		procPostMessageW.Call(uintptr(trayWnd), wmCommand, idQuit, 0)
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

func main() {
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
		appPort = port
		appWork = filepath.Join(os.Getenv("LOCALAPPDATA"), "TurnipProphet")
		openBrowser()
		return
	}
	mutex, already := acquireMutex()
	if mutex != 0 {
		defer windows.CloseHandle(mutex)
	}
	if already {
		for i := 0; i < 30; i++ {
			if port := FindRunning(); port > 0 {
				appPort = port
				appWork = filepath.Join(os.Getenv("LOCALAPPDATA"), "TurnipProphet")
				openBrowser()
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
	wc := wndClassEx{
		WndProc:   wndProcCB,
		Instance:  windows.Handle(hinst),
		ClassName: class,
	}
	wc.Size = uint32(unsafe.Sizeof(wc))
	if atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); atom == 0 {
		showError("Could not create the tray window. " + err.Error())
		return
	}
	hwnd, _, err := procCreateWindowExW.Call(0, uintptr(unsafe.Pointer(class)), uintptr(unsafe.Pointer(utf16Ptr("Turnip Prophet"))), 0, 0, 0, 0, 0, 0, 0, hinst, 0)
	if hwnd == 0 {
		showError("Could not create the tray window. " + err.Error())
		return
	}
	trayWnd = windows.HWND(hwnd)
	addTray(loadIcon())
	defer removeTray()
	openBrowser()
	runLoop()
	app.logf("Stopped.")
}
