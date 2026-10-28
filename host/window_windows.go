//go:build windows

package main

import (
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	wmSetIcon = 0x0080
	iconSmall = 0
	iconBig   = 1
	vtLPWSTR  = 31
)

var (
	procEnumWindows                 = user32.NewProc("EnumWindows")
	procGetWindowTextW              = user32.NewProc("GetWindowTextW")
	procGetClassNameW               = user32.NewProc("GetClassNameW")
	procIsWindowVisible             = user32.NewProc("IsWindowVisible")
	procSendMessageW                = user32.NewProc("SendMessageW")
	procSHGetPropertyStoreForWindow = shell32.NewProc("SHGetPropertyStoreForWindow")
	procDwmSetWindowAttribute       = windows.NewLazySystemDLL("dwmapi.dll").NewProc("DwmSetWindowAttribute")
	procRegGetValueW                = windows.NewLazySystemDLL("advapi32.dll").NewProc("RegGetValueW")
	procCoInitializeEx              = windows.NewLazySystemDLL("ole32.dll").NewProc("CoInitializeEx")
	procCoUninitialize              = windows.NewLazySystemDLL("ole32.dll").NewProc("CoUninitialize")

	iidPropertyStore = windows.GUID{
		Data1: 0x886d8eeb,
		Data2: 0x8cf2,
		Data3: 0x4446,
		Data4: [8]byte{0x8d, 0x02, 0xcd, 0xba, 0x1d, 0xbd, 0xcf, 0x99},
	}
	appUserModelFmt = windows.GUID{
		Data1: 0x9F4C2855,
		Data2: 0x9F79,
		Data3: 0x4B39,
		Data4: [8]byte{0xA8, 0xD0, 0xE1, 0xD4, 0x2D, 0xE1, 0xD5, 0xF3},
	}

	foundHWND uintptr
	enumCB    = syscall.NewCallback(enumAppWindow)
)

type propertyKey struct {
	Fmtid windows.GUID
	Pid   uint32
}

type propVariant struct {
	Vt        uint16
	Reserved1 uint16
	Reserved2 uint16
	Reserved3 uint16
	Val       uintptr
	Pad       uintptr
}

func enumAppWindow(hwnd, _ uintptr) uintptr {
	if vis, _, _ := procIsWindowVisible.Call(hwnd); vis == 0 {
		return 1
	}
	cls := make([]uint16, 64)
	procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&cls[0])), uintptr(len(cls)))
	if windows.UTF16ToString(cls) != "Chrome_WidgetWin_1" {
		return 1
	}
	buf := make([]uint16, 256)
	procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	title := windows.UTF16ToString(buf)
	if strings.HasPrefix(title, "Turnip Prophet") {
		foundHWND = hwnd
		return 0
	}
	return 1
}

func findAppWindow() windows.HWND {
	foundHWND = 0
	procEnumWindows.Call(enumCB, 0)
	return windows.HWND(foundHWND)
}

func setPropStr(store, setValue uintptr, pid uint32, text string) {
	ptr := utf16Ptr(text)
	if ptr == nil {
		return
	}
	key := propertyKey{Fmtid: appUserModelFmt, Pid: pid}
	pv := propVariant{Vt: vtLPWSTR, Val: uintptr(unsafe.Pointer(ptr))}
	syscall.SyscallN(setValue, store, uintptr(unsafe.Pointer(&key)), uintptr(unsafe.Pointer(&pv)))
	runtime.KeepAlive(ptr)
}

func stampWindow(hwnd windows.HWND) {
	exe, err := os.Executable()
	if err != nil || exe == "" {
		return
	}
	var store uintptr
	hr, _, _ := procSHGetPropertyStoreForWindow.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&iidPropertyStore)),
		uintptr(unsafe.Pointer(&store)),
	)
	if int32(hr) < 0 || store == 0 {
		return
	}
	vtbl := *(*uintptr)(unsafe.Pointer(store))
	methods := unsafe.Slice((*uintptr)(unsafe.Pointer(vtbl)), 8)
	setPropStr(store, methods[6], 5, "foctober.TurnipProphet")
	setPropStr(store, methods[6], 2, `"`+exe+`"`)
	setPropStr(store, methods[6], 3, exe+",0")
	setPropStr(store, methods[6], 4, "Turnip Prophet")
	syscall.SyscallN(methods[7], store)
	syscall.SyscallN(methods[2], store)
	setTitleBarTheme(hwnd)

	icon := loadIcon()
	if icon != 0 {
		procSendMessageW.Call(uintptr(hwnd), wmSetIcon, iconBig, uintptr(icon))
		procSendMessageW.Call(uintptr(hwnd), wmSetIcon, iconSmall, uintptr(icon))
	}
}

func appsPreferLight() bool {
	var val uint32 = 1
	size := uint32(unsafe.Sizeof(val))
	r, _, _ := procRegGetValueW.Call(
		0x80000001,
		uintptr(unsafe.Pointer(utf16Ptr(`Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`))),
		uintptr(unsafe.Pointer(utf16Ptr("AppsUseLightTheme"))),
		0x10,
		0,
		uintptr(unsafe.Pointer(&val)),
		uintptr(unsafe.Pointer(&size)),
	)
	if r != 0 {
		return true
	}
	return val != 0
}

func setTitleBarTheme(hwnd windows.HWND) {
	var dark int32
	if !appsPreferLight() {
		dark = 1
	}
	for _, attr := range []uintptr{20, 19} {
		procDwmSetWindowAttribute.Call(
			uintptr(hwnd),
			attr,
			uintptr(unsafe.Pointer(&dark)),
			unsafe.Sizeof(dark),
		)
	}
}

func watchAppWindow() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	hr, _, _ := procCoInitializeEx.Call(0, 2)
	if int32(hr) == 0 {
		defer procCoUninitialize.Call()
	}
	var hwnd windows.HWND
	for i := 0; i < 40; i++ {
		hwnd = findAppWindow()
		if hwnd != 0 {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}
	if hwnd == 0 {
		return
	}
	stampWindow(hwnd)
	time.Sleep(time.Second)
	if again := findAppWindow(); again != 0 {
		stampWindow(again)
	}
}
