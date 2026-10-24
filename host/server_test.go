package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStateRoundTripAndMigrate(t *testing.T) {
	root := t.TempDir()
	appData := t.TempDir()
	t.Setenv("APPDATA", appData)
	old := []byte(`{"weeks":[{"id":"2026-10-18","prices":[96]}],"activeId":"2026-10-18"}`)
	if err := os.WriteFile(filepath.Join(root, "turnips.json"), old, 0644); err != nil {
		t.Fatal(err)
	}
	app, err := Prepare(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "turnips.json.migrated")); err != nil {
		t.Fatal("old file was not renamed")
	}
	got, err := os.ReadFile(app.Data)
	if err != nil || string(got) != string(old) {
		t.Fatalf("copied data = %q %v", got, err)
	}
	srv := httptest.NewServer(app.Handler())
	defer srv.Close()
	res, err := http.Get(srv.URL + "/api/info")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if !strings.Contains(string(b), `"app":"turnip-prophet"`) || !strings.Contains(string(b), "turnips.json") {
		t.Fatal(string(b))
	}
	body := `{"weeks":[],"activeId":null}`
	res, err = http.Post(srv.URL+"/api/state", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatal(res.Status)
	}
	saved, _ := os.ReadFile(app.Data)
	if string(saved) != body {
		t.Fatal(string(saved))
	}
	res, err = http.Post(srv.URL+"/api/state", "application/json", strings.NewReader(`{"nope":1}`))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != 400 {
		t.Fatal(res.Status)
	}
}

func TestMissingFile(t *testing.T) {
	root := t.TempDir()
	miss := MissingFiles(root)
	if len(miss) != len(neededFiles) {
		t.Fatal(miss)
	}
}
