package app

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"dslparser/internal/cache"
	"dslparser/internal/discovery"
)

// stubHeaderReader представляет неиспользуемую в оркестраторных тестах границу чтения заголовков.
type stubHeaderReader struct{}

// ReadHeader возвращает пустой заголовок.
func (stubHeaderReader) ReadHeader(string) (cache.Header, error) { return cache.Header{}, nil }

// TestRunWithDependenciesProcessesDiscoveredFiles проверяет порядок, параметры и итог пакетной обработки.
func TestRunWithDependenciesProcessesDiscoveredFiles(t *testing.T) {
	scanErr := errors.New("scan failed")
	reader := stubHeaderReader{}
	wantTime := time.Date(2025, 2, 3, 4, 5, 6, 0, time.FixedZone("test", 3*60*60))
	paths := []string{"/input/a.txt", "/input/b.txt", "/input/c.txt"}
	jsonPaths := []string{"/input/a.json"}
	var requests []processRequest
	var gotDiscovery discovery.Request
	var gotCachePaths []string
	deps := dependencies{
		Discover: func(request discovery.Request) (discovery.Result, error) {
			gotDiscovery = request
			return discovery.Result{Files: paths, JSONFiles: jsonPaths, ScanErrors: []error{scanErr}}, nil
		},
		BuildCache: func(paths []string, gotReader cache.HeaderReader) (cache.Index, error) {
			gotCachePaths = append([]string(nil), paths...)
			if gotReader != reader {
				t.Fatalf("BuildCache reader = %#v, want supplied reader", gotReader)
			}
			return cache.Index{}, nil
		},
		HeaderReader: reader,
		Clock:        func() time.Time { return wantTime },
		processor: func(request processRequest) FileResult {
			requests = append(requests, request)
			switch request.Path {
			case paths[0]:
				return FileResult{Path: request.Path, Status: FileSuccess, Action: ActionCreated, ErrorCount: 2}
			case paths[1]:
				return FileResult{Path: request.Path, Status: FileSuccess, Action: ActionSkipped}
			default:
				return FileResult{Path: request.Path, Status: FileFailed, Action: ActionFailed, ErrorCount: 1}
			}
		},
	}
	depth := 2
	result, err := runWithDependencies(Options{Path: "input", Replace: true, Depth: &depth, ToolVersion: "v-test"}, deps)
	if err != nil {
		t.Fatalf("runWithDependencies() error: %v", err)
	}
	if gotDiscovery.Path != "input" || gotDiscovery.Depth != &depth {
		t.Fatalf("Discover request = %#v", gotDiscovery)
	}
	if !reflect.DeepEqual(gotCachePaths, jsonPaths) {
		t.Fatalf("BuildCache paths = %v, want %v", gotCachePaths, jsonPaths)
	}
	if len(requests) != len(paths) {
		t.Fatalf("processor calls = %d, want %d", len(requests), len(paths))
	}
	for index, request := range requests {
		if request.Path != paths[index] || !request.Replace || request.ToolVersion != "v-test" || !request.StartedAt.Equal(wantTime.UTC()) {
			t.Errorf("request[%d] = %#v", index, request)
		}
	}
	wantSummary := Summary{Found: 3, Parsed: 2, Created: 1, Skipped: 1, Success: 2, Failed: 1, Diagnostics: 3, ScanErrors: 1}
	if result.Summary != wantSummary || result.ExitCode != 1 {
		t.Fatalf("result summary/code = %#v/%d, want %#v/1", result.Summary, result.ExitCode, wantSummary)
	}
	if !reflect.DeepEqual(result.ScanErrors, []error{scanErr}) {
		t.Fatalf("ScanErrors = %v", result.ScanErrors)
	}
	result.ScanErrors[0] = errors.New("changed")
	if scanErr.Error() != "scan failed" {
		t.Fatal("RunResult unexpectedly aliases discovery scan-error slice")
	}
}

// TestRunWithDependenciesStopsBeforeProcessing проверяет ошибки запуска и построения кэша.
func TestRunWithDependenciesStopsBeforeProcessing(t *testing.T) {
	discoveryErr := errors.New("discovery failed")
	cacheErr := errors.New("cache failed")
	tests := []struct {
		name        string
		options     Options
		discoverErr error
		cacheErr    error
	}{
		{name: "invalid options", options: Options{ToolVersion: "v"}},
		{name: "discovery", options: Options{Path: "x", ToolVersion: "v"}, discoverErr: discoveryErr},
		{name: "cache", options: Options{Path: "x", ToolVersion: "v"}, cacheErr: cacheErr},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			processed := false
			_, err := runWithDependencies(test.options, dependencies{
				Discover: func(discovery.Request) (discovery.Result, error) {
					return discovery.Result{Files: []string{"x.txt"}}, test.discoverErr
				},
				BuildCache:   func([]string, cache.HeaderReader) (cache.Index, error) { return cache.Index{}, test.cacheErr },
				HeaderReader: stubHeaderReader{}, Clock: time.Now,
				processor: func(processRequest) FileResult { processed = true; return FileResult{} },
			})
			if err == nil {
				t.Fatal("runWithDependencies() error = nil")
			}
			if processed {
				t.Fatal("processor called after setup error")
			}
		})
	}
}

// TestRunWithDependenciesEmptyDiscovery проверяет успешный пустой пакет.
func TestRunWithDependenciesEmptyDiscovery(t *testing.T) {
	result, err := runWithDependencies(Options{Path: "x", ToolVersion: "v"}, dependencies{
		Discover:     func(discovery.Request) (discovery.Result, error) { return discovery.Result{}, nil },
		BuildCache:   func([]string, cache.HeaderReader) (cache.Index, error) { return cache.Index{}, nil },
		HeaderReader: stubHeaderReader{}, Clock: time.Now,
		processor: func(processRequest) FileResult { t.Fatal("unexpected processor call"); return FileResult{} },
	})
	if err != nil {
		t.Fatalf("runWithDependencies() error: %v", err)
	}
	if result.ExitCode != 0 || result.Summary != (Summary{}) || len(result.Files) != 0 {
		t.Fatalf("result = %#v, want empty successful result", result)
	}
}
