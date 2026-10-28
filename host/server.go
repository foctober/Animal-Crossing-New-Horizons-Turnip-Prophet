package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	portFirst = 47821
	portLast  = 47840
	maxBody   = 2_000_000
)

type fileSpec struct {
	rel  string
	kind string
}

var siteFiles = map[string]fileSpec{
	"/":               {filepath.Join("www", "index.html"), "text/html; charset=utf-8"},
	"/index.html":     {filepath.Join("www", "index.html"), "text/html; charset=utf-8"},
	"/styles.css":     {filepath.Join("www", "styles.css"), "text/css; charset=utf-8"},
	"/app.js":         {filepath.Join("www", "app.js"), "text/javascript; charset=utf-8"},
	"/predictions.js": {filepath.Join("predictions.js"), "text/javascript; charset=utf-8"},
	"/logo.png":       {filepath.Join("www", "logo.png"), "image/png"},
	"/field.jpg":      {filepath.Join("www", "field.jpg"), "image/jpeg"},
	"/wood.png":       {filepath.Join("www", "wood.png"), "image/png"},
}

var neededFiles = []string{
	"predictions.js",
	"TurnipProphet.ico",
	filepath.Join("www", "index.html"),
	filepath.Join("www", "app.js"),
	filepath.Join("www", "styles.css"),
	filepath.Join("www", "logo.png"),
	filepath.Join("www", "field.jpg"),
	filepath.Join("www", "wood.png"),
}

type App struct {
	Install     string
	Data        string
	Work        string
	logPath     string
	onQuit      func()
	onExitSaved func()
}

func MissingFiles(install string) []string {
	var miss []string
	for _, rel := range neededFiles {
		if _, err := os.Stat(filepath.Join(install, rel)); err != nil {
			miss = append(miss, rel)
		}
	}
	return miss
}

func (a *App) logf(format string, args ...any) {
	line := time.Now().Format("2006-01-02 15:04:05") + "  " + fmt.Sprintf(format, args...) + "\r\n"
	f, err := os.OpenFile(a.logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	_, _ = f.WriteString(line)
	_ = f.Close()
}

func Prepare(install string) (*App, error) {
	roaming := os.Getenv("APPDATA")
	if roaming == "" {
		roaming = install
	}
	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		local = roaming
	}
	dataDir := filepath.Join(roaming, "TurnipProphet")
	workDir := filepath.Join(local, "TurnipProphet")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(workDir, 0700); err != nil {
		return nil, err
	}
	dest := filepath.Join(dataDir, "turnips.json")
	app := &App{
		Install: install,
		Data:    dest,
		Work:    workDir,
		logPath: filepath.Join(workDir, "turnips.log"),
	}
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		for _, old := range legacyWeekFiles(install) {
			b, err := os.ReadFile(old)
			if err != nil || !validState(b) {
				continue
			}
			if err := os.WriteFile(dest, b, 0644); err != nil {
				return nil, err
			}
			bak := old + ".migrated"
			if _, err := os.Stat(bak); errors.Is(err, os.ErrNotExist) {
				_ = os.Rename(old, bak)
			}
			app.logf("Copied weeks to %s", dest)
			break
		}
	}
	return app, nil
}

func legacyWeekFiles(install string) []string {
	list := []string{filepath.Join(install, "turnips.json")}
	const previous = `C:\Tools\Turnips\turnips.json`
	if !strings.EqualFold(filepath.Clean(install), filepath.Clean(`C:\Tools\Turnips`)) {
		list = append(list, previous)
	}
	return list
}

func validState(b []byte) bool {
	if len(b) == 0 || len(b) > maxBody {
		return false
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(b, &raw) != nil {
		return false
	}
	_, ok := raw["weeks"]
	return ok
}

func (a *App) readState() []byte {
	b, err := os.ReadFile(a.Data)
	if err != nil || !validState(b) {
		return []byte(`{"weeks":[],"activeId":null}`)
	}
	return b
}

func replaceFile(tmp, dest string) error {
	if err := os.Rename(tmp, dest); err == nil {
		return nil
	}
	if err := os.Remove(dest); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.Rename(tmp, dest)
}

func (a *App) writeState(b []byte) error {
	tmp := a.Data + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return replaceFile(tmp, a.Data)
}

func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/info", a.handleInfo)
	mux.HandleFunc("/api/state", a.handleState)
	mux.HandleFunc("/api/shutdown", a.handleShutdown)
	mux.HandleFunc("/api/exit-saved", a.handleExitSaved)
	mux.HandleFunc("/favicon.ico", a.handleIcon)
	mux.HandleFunc("/TurnipProphet.ico", a.handleIcon)
	mux.HandleFunc("/", a.handleFile)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, b []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func (a *App) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(map[string]any{
		"app":     "turnip-prophet",
		"version": 1,
		"root":    a.Install,
		"data":    a.Data,
		"dark":    windowsAppDark(),
	})
	writeJSON(w, http.StatusOK, b)
}

func (a *App) handleState(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, a.readState())
	case http.MethodPost:
		b, err := io.ReadAll(io.LimitReader(r.Body, maxBody+1))
		if err != nil || !validState(b) {
			writeJSON(w, http.StatusBadRequest, []byte(`{"ok":false}`))
			return
		}
		if err := a.writeState(b); err != nil {
			a.logf("save failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, []byte(`{"ok":false}`))
			return
		}
		writeJSON(w, http.StatusOK, []byte(`{"ok":true}`))
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (a *App) handleExitSaved(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, []byte(`{"ok":true}`))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	if a.onExitSaved != nil {
		go a.onExitSaved()
	}
}

func (a *App) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, []byte(`{"ok":true}`))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	if a.onQuit != nil {
		go func() {
			time.Sleep(200 * time.Millisecond)
			a.onQuit()
		}()
	}
}

func (a *App) handleIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b, err := os.ReadFile(filepath.Join(a.Install, "TurnipProphet.ico"))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(b)
}

func (a *App) handleFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	spec, ok := siteFiles[r.URL.Path]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	b, err := os.ReadFile(filepath.Join(a.Install, spec.rel))
	if err != nil {
		http.Error(w, "missing", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", spec.kind)
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(b)
}

func FindRunning() int {
	client := &http.Client{Timeout: 400 * time.Millisecond}
	for p := portFirst; p <= portLast; p++ {
		res, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/api/info", p))
		if err != nil {
			continue
		}
		b, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
		res.Body.Close()
		if strings.Contains(string(b), `"app":"turnip-prophet"`) {
			return p
		}
	}
	return 0
}

func Listen() (net.Listener, int, error) {
	for p := portFirst; p <= portLast; p++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err == nil {
			return ln, p, nil
		}
	}
	return nil, 0, errors.New("no free local port between 47821 and 47840")
}

func writePortFile(work string, port int) {
	_ = os.WriteFile(filepath.Join(work, "turnips.port"), []byte(strconvItoa(port)), 0644)
}

func clearPortFile(work string) {
	_ = os.Remove(filepath.Join(work, "turnips.port"))
}

func strconvItoa(n int) string {
	return fmt.Sprintf("%d", n)
}
