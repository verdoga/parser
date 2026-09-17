package diagnostics

import (
	"errors"
	"reflect"
	"testing"

	"dslparser/internal/model"
)

// TestBuilderValidation проверяет отклонение всех нарушений контракта накопителя и запроса.
func TestBuilderValidation(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Fatal("New(\"\") error = nil")
	} else {
		assertRequestError(t, err, "source")
	}

	builder, err := New("processing-1")
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	location := testRange(1, 1, 1, 4)
	tests := []struct {
		name    string
		request Request
		field   string
	}{
		{name: "неизвестный код", request: Request{Code: Code("P999"), Scope: model.ScopeLine, Location: &location}, field: "code"},
		{name: "запрещённая область", request: Request{Code: P003, Scope: model.ScopeDocument, Location: nil, Details: Details{Fragment: "@bad"}}, field: "scope"},
		{name: "неизвестная область", request: Request{Code: P003, Scope: model.DiagnosticScope("unknown"), Location: &location, Details: Details{Fragment: "@bad"}}, field: "scope"},
		{name: "обязательный fatal", request: Request{Code: P001, Scope: model.ScopeDocument}, field: "fatal"},
		{name: "позиционная ошибка без диапазона", request: Request{Code: P003, Scope: model.ScopeElement, Details: Details{Fragment: "@bad"}}, field: "location"},
		{name: "документ с диапазоном", request: Request{Code: P001, Scope: model.ScopeDocument, Fatal: true, Location: &location}, field: "location"},
		{name: "нулевая строка", request: Request{Code: P003, Scope: model.ScopeElement, Location: rangePointer(testRange(0, 1, 1, 2)), Details: Details{Fragment: "@bad"}}, field: "location"},
		{name: "нулевая колонка", request: Request{Code: P003, Scope: model.ScopeElement, Location: rangePointer(testRange(1, 0, 1, 2)), Details: Details{Fragment: "@bad"}}, field: "location"},
		{name: "обратный диапазон", request: Request{Code: P003, Scope: model.ScopeElement, Location: rangePointer(testRange(2, 1, 1, 2)), Details: Details{Fragment: "@bad"}}, field: "location"},
		{name: "неверный дополнительный диапазон", request: Request{Code: P011, Scope: model.ScopeBlock, Fatal: true, Location: &location, RelatedLocations: []model.Range{testRange(0, 1, 1, 2)}, Details: Details{Tag: "text"}}, field: "relatedLocations"},
		{name: "дополнительный равен основному", request: Request{Code: P011, Scope: model.ScopeBlock, Fatal: true, Location: &location, RelatedLocations: []model.Range{location}, Details: Details{Tag: "text"}}, field: "relatedLocations"},
		{name: "повтор дополнительных", request: Request{Code: P011, Scope: model.ScopeBlock, Fatal: true, Location: &location, RelatedLocations: []model.Range{testRange(2, 1, 2, 2), testRange(2, 1, 2, 2)}, Details: Details{Tag: "text"}}, field: "relatedLocations"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := builder.Add(test.request); err == nil {
				t.Fatal("Add() error = nil")
			} else {
				assertRequestError(t, err, test.field)
			}
		})
	}

	detailTests := []struct {
		name    string
		code    Code
		details Details
	}{
		{name: "P003 fragment", code: P003},
		{name: "P004 required", code: P004},
		{name: "P005 tag", code: P005, details: Details{Received: "block"}},
		{name: "P005 received", code: P005, details: Details{Tag: "question"}},
		{name: "P006 tag", code: P006, details: Details{Required: "ID"}},
		{name: "P006 required", code: P006, details: Details{Tag: "task"}},
		{name: "P007 fragment", code: P007},
		{name: "P008 fragment", code: P008},
		{name: "P010 fragment", code: P010},
		{name: "P011 tag", code: P011},
		{name: "P012 fragment", code: P012},
		{name: "P014 received", code: P014, details: Details{Supported: []string{"1.2"}}},
		{name: "P014 supported nil", code: P014, details: Details{Received: "2.0"}},
		{name: "P014 supported empty value", code: P014, details: Details{Received: "2.0", Supported: []string{""}}},
		{name: "IO001 path", code: IO001, details: Details{Cause: errors.New("denied")}},
		{name: "IO001 cause", code: IO001, details: Details{Path: "/lesson.txt"}},
	}
	for _, test := range detailTests {
		t.Run(test.name, func(t *testing.T) {
			if err := validateDetails(test.code, test.details); err == nil {
				t.Fatal("validateDetails() error = nil")
			} else {
				assertRequestError(t, err, "details")
			}
		})
	}
}

// TestRangeHelpers проверяет сравнение и глубокое копирование диапазонов диагностик.
func TestRangeHelpers(t *testing.T) {
	first := testRange(1, 2, 3, 4)
	equal := first
	different := testRange(1, 2, 3, 5)
	if !sameLocation(nil, nil) || !sameLocation(&first, &equal) {
		t.Fatal("sameLocation() rejected equal locations")
	}
	if sameLocation(&first, nil) || sameLocation(nil, &first) || sameLocation(&first, &different) {
		t.Fatal("sameLocation() accepted different locations")
	}
	if cloneRange(nil) != nil {
		t.Fatal("cloneRange(nil) is not nil")
	}
	cloned := cloneRange(&first)
	cloned.Start.Line = 99
	if first.Start.Line != 1 {
		t.Fatal("cloneRange() aliases its input")
	}
	if clonedRanges := cloneRanges(nil); clonedRanges == nil || len(clonedRanges) != 0 {
		t.Fatalf("cloneRanges(nil) = %#v, want non-nil empty slice", clonedRanges)
	}
	ranges := []model.Range{first}
	clonedRanges := cloneRanges(ranges)
	clonedRanges[0].Start.Line = 99
	if ranges[0].Start.Line != 1 {
		t.Fatal("cloneRanges() aliases its input")
	}
}

// TestBuilderOrderAndDeduplication проверяет порядок, ID и ключ дедупликации code/location.
func TestBuilderOrderAndDeduplication(t *testing.T) {
	builder, err := New("processing-7")
	if err != nil {
		t.Fatal(err)
	}
	firstLocation := testRange(2, 1, 2, 9)
	secondLocation := testRange(3, 1, 3, 2)
	requests := []builderCase{
		{name: "первая проблема", request: Request{Code: P003, Scope: model.ScopeElement, Location: &firstLocation, Details: Details{Fragment: "@unknown"}}, created: true},
		{name: "дубликат с другими агрегатами", request: Request{Code: P003, Scope: model.ScopeElement, Fatal: true, Location: rangePointer(firstLocation), RelatedLocations: []model.Range{secondLocation}, Details: Details{Fragment: "другой текст"}}, created: false},
		{name: "тот же диапазон и другой код", request: Request{Code: P012, Scope: model.ScopeElement, Location: rangePointer(firstLocation), Details: Details{Fragment: "{"}}, created: true},
		{name: "тот же код и другой диапазон", request: Request{Code: P003, Scope: model.ScopeElement, Location: &secondLocation, Details: Details{Fragment: "@other"}}, created: true},
	}
	var duplicate model.Diagnostic
	for _, test := range requests {
		t.Run(test.name, func(t *testing.T) {
			got, created, err := builder.Add(test.request)
			if err != nil {
				t.Fatalf("Add() error: %v", err)
			}
			if created != test.created {
				t.Fatalf("Add() created = %v, want %v", created, test.created)
			}
			if !created {
				duplicate = got
			}
		})
	}
	items := builder.Items()
	if len(items) != 3 {
		t.Fatalf("Items() length = %d, want 3", len(items))
	}
	if items[0].Code != string(P003) || items[1].Code != string(P012) || items[2].Code != string(P003) {
		t.Fatalf("Items() codes = %q, %q, %q", items[0].Code, items[1].Code, items[2].Code)
	}
	for position, item := range items {
		wantID := diagnosticID("processing-7", position+1)
		if item.ID != wantID || item.Source != "processing-7" || item.Severity != model.SeverityError {
			t.Errorf("Items()[%d] = %#v, want ID %q and error source processing-7", position, item, wantID)
		}
	}
	if duplicate.ID != items[0].ID || duplicate.Message != items[0].Message || duplicate.Fatal != items[0].Fatal || !reflect.DeepEqual(duplicate.RelatedLocations, items[0].RelatedLocations) {
		t.Fatalf("duplicate = %#v, want copy of %#v", duplicate, items[0])
	}
	if diagnosticID("processing-7", 1) == diagnosticID("processing-8", 1) || diagnosticID("processing-7", 1) == diagnosticID("processing-7", 2) {
		t.Fatal("diagnosticID() is not unique by source and ordinal")
	}
}

// TestBuilderOwnership проверяет отсутствие общих указателей и срезов с запросом и результатами.
func TestBuilderOwnership(t *testing.T) {
	builder, err := New("processing-owned")
	if err != nil {
		t.Fatal(err)
	}
	location := testRange(1, 1, 1, 4)
	related := []model.Range{testRange(2, 1, 2, 3)}
	supported := []string{"1.2", "1.3"}
	request := Request{
		Code:             P014,
		Scope:            model.ScopeElement,
		Fatal:            true,
		Location:         &location,
		RelatedLocations: related,
		Details:          Details{Received: "2.0", Supported: supported},
	}
	diagnostic, created, err := builder.Add(request)
	if err != nil || !created {
		t.Fatalf("Add() = %#v, %v, %v", diagnostic, created, err)
	}
	wantMessage := diagnostic.Message
	location.Start.Line = 99
	related[0].Start.Line = 99
	supported[0] = "changed"
	items := builder.Items()
	if items[0].Location.Start.Line != 1 || items[0].RelatedLocations[0].Start.Line != 2 || items[0].Message != wantMessage {
		t.Fatalf("request mutation changed stored diagnostic: %#v", items[0])
	}

	diagnostic.Location.Start.Line = 77
	diagnostic.RelatedLocations[0].Start.Line = 77
	items[0].Location.Start.Line = 88
	items[0].RelatedLocations[0].Start.Line = 88
	items = append(items, model.Diagnostic{})
	again := builder.Items()
	if len(again) != 1 || again[0].Location.Start.Line != 1 || again[0].RelatedLocations[0].Start.Line != 2 {
		t.Fatalf("returned diagnostic aliases builder state: %#v", again)
	}
	if again[0].RelatedLocations == nil {
		t.Fatal("Items() returned nil required relatedLocations slice")
	}
}

// assertRequestError проверяет устойчивый тип и поле ошибки запроса.
func assertRequestError(t *testing.T, err error, field string) {
	t.Helper()
	var requestError *RequestError
	if !errors.As(err, &requestError) {
		t.Fatalf("error type = %T, want *RequestError", err)
	}
	if requestError.Field != field || requestError.Rule == "" {
		t.Fatalf("RequestError = %#v, want field %q and non-empty rule", requestError, field)
	}
}

// testRange создаёт диапазон для тестов накопителя.
func testRange(startLine, startColumn, endLine, endColumn int) model.Range {
	return model.Range{Start: model.Position{Line: startLine, Column: startColumn}, End: model.Position{Line: endLine, Column: endColumn}}
}

// rangePointer возвращает указатель на отдельную копию диапазона.
func rangePointer(value model.Range) *model.Range { return &value }
