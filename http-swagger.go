package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
	"github.com/fsnotify/fsnotify"
	"github.com/swaggo/http-swagger"
)

type SwaggerLoader struct {
	mu        sync.RWMutex
	data      []byte
	source    string
	isRemote  bool
	pollInterval time.Duration
}

func NewSwaggerLoader(source string, pollInterval time.Duration) *SwaggerLoader {
	sl := &SwaggerLoader{
		source:       source,
		isRemote:     isURL(source),
		pollInterval: pollInterval,
	}
	if err := sl.Reload(); err != nil {
		log.Fatalf("Initial load failed: %v", err)
	}
	sl.StartWatcher()
	return sl
}

func (sl *SwaggerLoader) Reload() error {
	var data []byte
	var err error

	if sl.isRemote {
		data, err = sl.fetchRemote()
	} else {
		data, err = os.ReadFile(sl.source)
	}

	if err != nil {
		return fmt.Errorf("load failed: %w", err)
	}

	if !json.Valid(data) {
		return fmt.Errorf("invalid JSON format")
	}

	sl.mu.Lock()
	sl.data = data
	sl.mu.Unlock()
	return nil
}

func (sl *SwaggerLoader) fetchRemote() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", sl.source, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (sl *SwaggerLoader) StartWatcher() {
	if sl.isRemote {
		go sl.remotePolling()
	} else {
		go sl.localFileWatch()
	}
}

func (sl *SwaggerLoader) remotePolling() {
	ticker := time.NewTicker(sl.pollInterval)
	for range ticker.C {
		if err := sl.Reload(); err != nil {
			log.Printf("Remote reload error: %v", err)
		}
	}
}

func (sl *SwaggerLoader) localFileWatch() {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	watcher.Add(sl.source)
	
	var debounceTimer *time.Timer
	for {
		select {
		case event := <-watcher.Events:
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Rename) {
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(500*time.Millisecond, func() {
					if err := sl.Reload(); err != nil {
						log.Printf("Local reload error: %v", err)
					}
				})
			}
		case err := <-watcher.Errors:
			log.Printf("Watcher error: %v", err)
		}
	}
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func main() {
	var source string
	flag.StringVar(&source, "source", "./swagger.json", "JSON源路径/URL")
	flag.Parse()

	// 自动检测当前目录
	if !isURL(source) {
		if exe, err := os.Executable(); err == nil {
			source = filepath.Join(filepath.Dir(exe), source)
		}
	}

	loader := NewSwaggerLoader(source, 30*time.Second)

	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		loader.mu.RLock()
		defer loader.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		w.Write(loader.data)
	})

	http.Handle("/docs/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
		httpSwagger.UIConfig(map[string]string{
			"persistAuthorization": "true",
		}),
	))

	log.Println("Server started on :8585")
	log.Fatal(http.ListenAndServe(":8585", nil))
}
