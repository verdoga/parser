package discovery

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

// TestDiscoverPathKinds проверяет все допустимые и недопустимые виды исходного пути.
func TestDiscoverPathKinds(t *testing.T) {
	root := t.TempDir()
	txt := filepath.Join(root, "lesson.TxT")
	writeFile(t, txt)
	directory := filepath.Join(root, "lessons")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	nonTXT := filepath.Join(root, "lesson.md")
	writeFile(t, nonTXT)
	link := filepath.Join(root, "lesson-link.txt")
	if err := os.Symlink(txt, link); err != nil {
		t.Fatal(err)
	}
	fifo := filepath.Join(root, "lesson.txt.fifo")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		path       string
		wantRoot   string
		wantIsFile bool
		wantFiles  []string
		wantError  bool
	}{
		{name: "regular mixed-case TXT", path: txt, wantRoot: txt, wantIsFile: true, wantFiles: []string{txt}},
		{name: "directory", path: directory, wantRoot: directory, wantFiles: []string{}},
		{name: "missing", path: filepath.Join(root, "missing.txt"), wantError: true},
		{name: "regular non-TXT", path: nonTXT, wantError: true},
		{name: "symbolic link", path: link, wantError: true},
		{name: "special file", path: fifo, wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Discover(Request{Path: test.path})
			if test.wantError {
				assertInputError(t, err)
				return
			}
			if err != nil {
				t.Fatalf("Discover() error: %v", err)
			}
			if result.Root != test.wantRoot || result.IsFile != test.wantIsFile || !reflect.DeepEqual(result.Files, test.wantFiles) {
				t.Fatalf("Discover() = %#v, want root %q, IsFile %t, files %q", result, test.wantRoot, test.wantIsFile, test.wantFiles)
			}
			if result.Files == nil || result.ScanErrors == nil && !test.wantIsFile {
				t.Fatalf("Discover() returned nil required slices: %#v", result)
			}
		})
	}

	_, err := Discover(Request{Path: filepath.Join(root, "missing.txt")})
	var pathError *os.PathError
	if !errors.As(err, &pathError) || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Discover() error = %v, want wrapped *os.PathError matching os.ErrNotExist", err)
	}
}

// TestDiscoverDepthAndOrder проверяет глубину, сортировку и игнорирование ссылок.
func TestDiscoverDepthAndOrder(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "z.TXT"))
	writeFile(t, filepath.Join(root, "a.txt"))
	writeFile(t, filepath.Join(root, "ignore.md"))
	writeFile(t, filepath.Join(root, "first", "b.txt"))
	writeFile(t, filepath.Join(root, "first", "second", "c.txt"))
	outside := filepath.Join(t.TempDir(), "outside.txt")
	writeFile(t, outside)
	if err := os.Symlink(outside, filepath.Join(root, "linked-file.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Dir(outside), filepath.Join(root, "linked-directory")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		depth *int
		want  []string
	}{
		{name: "unlimited", want: absolutePaths(root, "a.txt", filepath.Join("first", "b.txt"), filepath.Join("first", "second", "c.txt"), "z.TXT")},
		{name: "zero", depth: intPointer(0), want: absolutePaths(root, "a.txt", "z.TXT")},
		{name: "one", depth: intPointer(1), want: absolutePaths(root, "a.txt", filepath.Join("first", "b.txt"), "z.TXT")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Discover(Request{Path: filepath.Join(root, "."), Depth: test.depth})
			if err != nil {
				t.Fatalf("Discover() error: %v", err)
			}
			if result.Root != filepath.Clean(root) {
				t.Fatalf("Discover() root = %q, want %q", result.Root, filepath.Clean(root))
			}
			if !reflect.DeepEqual(result.Files, test.want) {
				t.Fatalf("Discover() files = %q, want %q", result.Files, test.want)
			}
		})
	}
}

// TestDiscoverScanErrors проверяет продолжение обхода и ошибку недоступного корня.
func TestDiscoverScanErrors(t *testing.T) {
	if os.Getenv("DISCOVERY_PERMISSION_HELPER") != "" {
		runPermissionHelper(t)
		return
	}

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a", "visible.txt"))
	writeFile(t, filepath.Join(root, "b-denied", "hidden.txt"))
	writeFile(t, filepath.Join(root, "c", "visible.txt"))
	if err := os.Chmod(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, directory := range []string{"a", "c"} {
		if err := os.Chmod(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(filepath.Join(root, "b-denied"), 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(filepath.Join(root, "b-denied"), 0o755) })

	helperBinary := filepath.Join(root, "discovery.test")
	binary, err := os.ReadFile(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(helperBinary, binary, 0o755); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(helperBinary, "-test.run=^TestDiscoverScanErrors$")
	command.Env = append(os.Environ(), "DISCOVERY_PERMISSION_HELPER=child", "DISCOVERY_PERMISSION_ROOT="+root)
	if os.Geteuid() == 0 {
		command.SysProcAttr = &syscall.SysProcAttr{Credential: &syscall.Credential{Uid: 65534, Gid: 65534}}
	}
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("permission helper failed: %v\n%s", err, output)
	}
}

// TestDiscoverJSONScope проверяет область, дедупликацию и сортировку JSON-файлов.
func TestDiscoverJSONScope(t *testing.T) {
	root := t.TempDir()
	first := filepath.Join(root, "first")
	second := filepath.Join(root, "second")
	foreign := filepath.Join(root, "foreign")
	txtFiles := []string{
		filepath.Join(first, "b.txt"),
		filepath.Join(first, "a.txt"),
		filepath.Join(second, "c.txt"),
	}
	for _, path := range txtFiles {
		writeFile(t, path)
	}
	for _, path := range []string{
		filepath.Join(first, "z.JSON"),
		filepath.Join(first, "a.json"),
		filepath.Join(second, "m.JsOn"),
		filepath.Join(foreign, "ignored.json"),
		filepath.Join(first, "ignored.json.bak"),
	} {
		writeFile(t, path)
	}
	if err := os.Symlink(filepath.Join(first, "a.json"), filepath.Join(first, "linked.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(first, "directory.json"), 0o755); err != nil {
		t.Fatal(err)
	}

	want := []string{filepath.Join(first, "a.json"), filepath.Join(first, "z.JSON"), filepath.Join(second, "m.JsOn")}
	files, scanErrors := discoverJSONFiles(txtFiles)
	if len(scanErrors) != 0 {
		t.Fatalf("discoverJSONFiles() errors = %v", scanErrors)
	}
	if !reflect.DeepEqual(files, want) {
		t.Fatalf("discoverJSONFiles() = %q, want %q", files, want)
	}

	directoryResult, err := Discover(Request{Path: root})
	if err != nil {
		t.Fatalf("Discover(directory) error: %v", err)
	}
	if !reflect.DeepEqual(directoryResult.JSONFiles, want) {
		t.Fatalf("Discover(directory) JSON files = %q, want %q", directoryResult.JSONFiles, want)
	}

	fileResult, err := Discover(Request{Path: txtFiles[0]})
	if err != nil {
		t.Fatalf("Discover(file) error: %v", err)
	}
	wantFileScope := want[:2]
	if !reflect.DeepEqual(fileResult.JSONFiles, wantFileScope) {
		t.Fatalf("Discover(file) JSON files = %q, want %q", fileResult.JSONFiles, wantFileScope)
	}
}

// runPermissionHelper выполняет проверки ошибок доступа с правами непривилегированного пользователя.
func runPermissionHelper(t *testing.T) {
	root := os.Getenv("DISCOVERY_PERMISSION_ROOT")
	result, err := Discover(Request{Path: root})
	if err != nil {
		t.Fatalf("Discover(readable root) error: %v", err)
	}
	wantFiles := absolutePaths(root, filepath.Join("a", "visible.txt"), filepath.Join("c", "visible.txt"))
	if !reflect.DeepEqual(result.Files, wantFiles) {
		t.Fatalf("Discover(readable root) files = %q, want %q", result.Files, wantFiles)
	}
	if len(result.ScanErrors) != 1 || !strings.Contains(result.ScanErrors[0].Error(), filepath.Join(root, "b-denied")) {
		t.Fatalf("Discover(readable root) scan errors = %v, want one ordered b-denied error", result.ScanErrors)
	}

	_, err = Discover(Request{Path: filepath.Join(root, "b-denied")})
	assertInputError(t, err)
}

// assertInputError проверяет тип ошибки исходного пути.
func assertInputError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Discover() error = nil, want *InputError")
	}
	var inputErr *InputError
	if !errors.As(err, &inputErr) {
		t.Fatalf("Discover() error type = %T, want *InputError", err)
	}
}

// writeFile создаёт пустой файл вместе с отсутствующими родительскими каталогами.
func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

// absolutePaths строит ожидаемые абсолютные пути относительно корневого каталога.
func absolutePaths(root string, paths ...string) []string {
	result := make([]string, len(paths))
	for index, path := range paths {
		result[index] = filepath.Join(root, path)
	}
	return result
}

// intPointer возвращает указатель на отдельную копию целого числа.
func intPointer(value int) *int { return &value }
