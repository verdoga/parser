package app

import "testing"

// TestFileResultStrings проверяет стабильные машинные представления перечислений.
func TestFileResultStrings(t *testing.T) {
	statusTests := []struct {
		status FileStatus
		want   string
	}{{FileSuccess, "success"}, {FileFailed, "failed"}}
	for _, test := range statusTests {
		t.Run(test.want, func(t *testing.T) {
			if got := test.status.String(); got != test.want {
				t.Errorf("FileStatus(%d).String() = %q, want %q", test.status, got, test.want)
			}
		})
	}
	actionTests := []struct {
		action FileAction
		want   string
	}{{ActionCreated, "created"}, {ActionReplaced, "replaced"}, {ActionSkipped, "skipped"}, {ActionFailed, "failed"}}
	for _, test := range actionTests {
		t.Run(test.want, func(t *testing.T) {
			if got := test.action.String(); got != test.want {
				t.Errorf("FileAction(%d).String() = %q, want %q", test.action, got, test.want)
			}
		})
	}
	if got := FileStatus(99).String(); got != "unknown" {
		t.Errorf("unknown FileStatus.String() = %q, want unknown", got)
	}
	if got := FileAction(99).String(); got != "unknown" {
		t.Errorf("unknown FileAction.String() = %q, want unknown", got)
	}
}

// TestSummary проверяет учёт всех действий, диагностик и ошибок обхода.
func TestSummary(t *testing.T) {
	var summary Summary
	results := []FileResult{
		{Status: FileSuccess, Action: ActionCreated, ErrorCount: 2},
		{Status: FileSuccess, Action: ActionReplaced, ErrorCount: 1},
		{Status: FileSuccess, Action: ActionSkipped},
		{Status: FileFailed, Action: ActionFailed, ErrorCount: 3},
	}
	for _, result := range results {
		summary.addFile(result)
	}
	summary.addScanErrors(2)
	want := Summary{Found: 4, Parsed: 3, Created: 1, Replaced: 1, Skipped: 1, Success: 3, Failed: 1, Diagnostics: 6, ScanErrors: 2}
	if summary != want {
		t.Fatalf("Summary = %#v, want %#v", summary, want)
	}
	if got := summary.code(); got != 1 {
		t.Fatalf("code() = %d, want 1", got)
	}
	if got := (Summary{Success: 1}).code(); got != 0 {
		t.Fatalf("successful code() = %d, want 0", got)
	}
	if got := (Summary{ScanErrors: 1}).code(); got != 1 {
		t.Fatalf("scan-error code() = %d, want 1", got)
	}
}

// TestOptionsValidate проверяет только значения параметров, не обращаясь к пути.
func TestOptionsValidate(t *testing.T) {
	negative := -1
	zero := 0
	tests := []struct {
		name    string
		options Options
		wantErr bool
	}{
		{name: "valid nonexistent path", options: Options{Path: "/definitely/not/read", ToolVersion: "test"}},
		{name: "zero depth", options: Options{Path: "x", Depth: &zero, ToolVersion: "test"}},
		{name: "empty path", options: Options{ToolVersion: "test"}, wantErr: true},
		{name: "negative depth", options: Options{Path: "x", Depth: &negative, ToolVersion: "test"}, wantErr: true},
		{name: "empty tool version", options: Options{Path: "x"}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.options.validate()
			if (err != nil) != test.wantErr {
				t.Fatalf("validate() error = %v, wantErr %t", err, test.wantErr)
			}
		})
	}
}
