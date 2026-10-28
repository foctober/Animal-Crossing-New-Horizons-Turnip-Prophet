//go:build windows

package main

import (
	"fmt"
	"path/filepath"

	"github.com/jchv/go-webview2/pkg/edge"
)

var (
	browser         *edge.Chromium
	procGetKeyState = user32.NewProc("GetKeyState")
)

func blockZoomKey(vk uint) bool {
	switch vk {
	case 0xBB, 0xBD, 0x6B, 0x6D, 0x30, 0x60:
	default:
		return false
	}
	state, _, _ := procGetKeyState.Call(0x11)
	return int16(state) < 0
}

func attachPage() {
	page := edge.NewChromium()
	page.DataPath = filepath.Join(appWork, "webview")
	page.AcceleratorKeyCallback = blockZoomKey
	if !page.Embed(uintptr(frameWnd)) {
		procShowWindow.Call(uintptr(frameWnd), swHide)
		openBrowser()
		return
	}
	browser = page
	if settings, err := browser.GetSettings(); err == nil && settings != nil {
		_ = settings.PutIsZoomControlEnabled(false)
		_ = settings.PutIsPinchZoomEnabled(false)
		_ = settings.PutAreBrowserAcceleratorKeysEnabled(false)
	}
	browser.Resize()
	browser.Navigate(fmt.Sprintf("http://127.0.0.1:%d/", appPort))
	browser.Focus()
}
