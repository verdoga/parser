package model

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// TestValidateDocument проверяет сведения об источнике, метаданные и вычисляемый признак ошибок.
func TestValidateDocument(t *testing.T) {
	valid := validResult()
	if err := validateDocument(valid); err != nil {
		t.Fatalf("validateDocument(valid) error: %v", err)
	}

	tests := []struct {
		name   string
		change func(*Result)
	}{
		{name: "версия формата", change: func(r *Result) { r.FormatVersion = "2.0" }},
		{name: "пустое имя файла", change: func(r *Result) { r.Document.FileName = stringPtr("") }},
		{name: "неабсолютный путь", change: func(r *Result) { r.Document.FilePath = stringPtr("lesson.txt") }},
		{name: "имя не совпадает с путём", change: func(r *Result) { r.Document.FileName = stringPtr("other.txt") }},
		{name: "неподдерживаемая кодировка", change: func(r *Result) { r.Document.Encoding = stringPtr("UTF-16") }},
		{name: "отрицательное число строк", change: func(r *Result) { r.Document.LineCount = intPtr(-1) }},
		{name: "отрицательный размер", change: func(r *Result) { r.Document.ByteLength = intPtr(-1) }},
		{name: "короткий sha256", change: func(r *Result) { r.Document.SHA256 = stringPtr("abcd") }},
		{name: "sha256 в верхнем регистре", change: func(r *Result) { r.Document.SHA256 = stringPtr(strings.Repeat("A", 64)) }},
		{name: "пустой document id", change: func(r *Result) { r.Document.Metadata.DocumentID = stringPtr("") }},
		{name: "пустой заголовок", change: func(r *Result) { r.Document.Metadata.Title = stringPtr("") }},
		{name: "пустой раздел", change: func(r *Result) { r.Document.Metadata.Section = stringPtr("") }},
		{name: "order равен нулю", change: func(r *Result) { r.Document.Metadata.Order = stringPtr("0") }},
		{name: "order содержит не цифры", change: func(r *Result) { r.Document.Metadata.Order = stringPtr("1a") }},
		{name: "пустой resource dir", change: func(r *Result) { r.Document.Metadata.ResourceDirs = []string{"assets", ""} }},
		{name: "лишний hasErrors", change: func(r *Result) { r.Document.HasErrors = true }},
		{name: "пропущенный hasErrors", change: func(r *Result) {
			r.Document.HasErrors = false
			r.Diagnostics = []Diagnostic{documentDiagnostic(SeverityError)}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validResult()
			test.change(&result)
			assertInvariantError(t, validateDocument(result))
		})
	}

	t.Run("недоступный источник", func(t *testing.T) {
		result := Result{FormatVersion: FormatVersion, Document: Document{Metadata: Metadata{ResourceDirs: nil}}}
		if err := validateDocument(result); err != nil {
			t.Fatalf("validateDocument(unavailable) error: %v", err)
		}
	})
}

// TestValidateProcessing проверяет идентификаторы, временные поля и порядок запусков.
func TestValidateProcessing(t *testing.T) {
	valid := []Processing{
		{ID: "run-1", Tool: "dsl-parser", Version: stringPtr("1.0.0"), StartedAt: stringPtr("2026-09-17T10:20:30Z"), DurationMS: intPtr(0)},
		{ID: "run-2", Tool: "dsl-validator", Version: nil, StartedAt: stringPtr("2026-09-17T13:20:30+03:00"), DurationMS: nil},
	}
	if err := validateProcessing(valid); err != nil {
		t.Fatalf("validateProcessing(valid) error: %v", err)
	}

	tests := []struct {
		name  string
		items []Processing
	}{
		{name: "пустой id", items: []Processing{{Tool: "tool"}}},
		{name: "повтор id", items: []Processing{{ID: "same", Tool: "one"}, {ID: "same", Tool: "two"}}},
		{name: "пустой tool", items: []Processing{{ID: "run"}}},
		{name: "пустая version", items: []Processing{{ID: "run", Tool: "tool", Version: stringPtr("")}}},
		{name: "не RFC3339", items: []Processing{{ID: "run", Tool: "tool", StartedAt: stringPtr("17.09.2026")}}},
		{name: "время без часового пояса", items: []Processing{{ID: "run", Tool: "tool", StartedAt: stringPtr("2026-09-17T10:20:30")}}},
		{name: "отрицательная длительность", items: []Processing{{ID: "run", Tool: "tool", DurationMS: intPtr(-1)}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) { assertInvariantError(t, validateProcessing(test.items)) })
	}
}

// TestValidateLines проверяет структуру строк, Unicode-диапазоны и правила элементов.
func TestValidateLines(t *testing.T) {
	document, lines := validLines()
	if err := validateLines(document, lines); err != nil {
		t.Fatalf("validateLines(valid) error: %v", err)
	}

	tests := []struct {
		name   string
		change func(*Document, *[]Line)
	}{
		{name: "lineCount не совпадает", change: func(d *Document, _ *[]Line) { d.LineCount = intPtr(2) }},
		{name: "номер не от единицы", change: func(_ *Document, l *[]Line) { (*l)[0].Number = 0 }},
		{name: "пропуск номера", change: func(_ *Document, l *[]Line) { (*l)[1].Number = 3 }},
		{name: "неизвестный тип строки", change: func(_ *Document, l *[]Line) { (*l)[0].Type = LineType("other") }},
		{name: "отрицательная вложенность", change: func(_ *Document, l *[]Line) { (*l)[0].NestingLevel = -1 }},
		{name: "родитель следует после потомка", change: func(_ *Document, l *[]Line) { (*l)[0].ParentLine = intPtr(2) }},
		{name: "несуществующий родитель", change: func(_ *Document, l *[]Line) { (*l)[1].ParentLine = intPtr(9) }},
		{name: "неверная вложенность потомка", change: func(_ *Document, l *[]Line) { (*l)[1].NestingLevel = 2 }},
		{name: "у heading есть родитель", change: func(_ *Document, l *[]Line) { (*l)[0].ParentLine = intPtr(1) }},
		{name: "у heading ненулевая вложенность", change: func(_ *Document, l *[]Line) { (*l)[0].NestingLevel = 1 }},
		{name: "неизвестный eol", change: func(_ *Document, l *[]Line) { (*l)[0].EOL = LineEnding("\r") }},
		{name: "eof до последней строки", change: func(_ *Document, l *[]Line) { (*l)[0].EOL = EOLNone }},
		{name: "неизвестный тип элемента", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[0].Type = ElementType("other") }},
		{name: "start равен нулю", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[0].Start = 0 }},
		{name: "end меньше start", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[0].End = 1 }},
		{name: "диапазон за строкой", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[0].End = 99 }},
		{name: "raw не совпадает", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[1].Raw = "другой" }},
		{name: "элементы не отсортированы", change: func(_ *Document, l *[]Line) {
			(*l)[0].Elements[0], (*l)[0].Elements[1] = (*l)[0].Elements[1], (*l)[0].Elements[0]
		}},
		{name: "элементы пересекаются", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[1].Start = 2 }},
		{name: "повтор error id", change: func(_ *Document, l *[]Line) { (*l)[0].Elements[0].ErrorIDs = []string{"d1", "d1"} }},
		{name: "invalid без unparsed", change: func(_ *Document, l *[]Line) { (*l)[1].Type = LineInvalid }},
		{name: "unparsed без ошибки", change: func(_ *Document, l *[]Line) {
			(*l)[1].Elements[0].Type = ElementUnparsed
			(*l)[1].Elements[0].Value = nil
		}},
		{name: "block open с неверным raw", change: func(_ *Document, l *[]Line) {
			(*l)[1].Elements[0] = Element{Type: ElementBlockOpen, Raw: "т", Value: stringPtr("{"), Start: 1, End: 2, ErrorIDs: []string{}}
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc, gotLines := validLines()
			test.change(&doc, &gotLines)
			assertInvariantError(t, validateLines(doc, gotLines))
		})
	}

	t.Run("пустой файл", func(t *testing.T) {
		zero := 0
		if err := validateLines(Document{LineCount: &zero}, []Line{}); err != nil {
			t.Fatalf("validateLines(empty) error: %v", err)
		}
	})
}

// TestValidateDiagnostics проверяет ссылки, диапазоны, порядок и вычисленные признаки ошибок.
func TestValidateDiagnostics(t *testing.T) {
	processing, lines, diagnostics := validDiagnostics()
	if err := validateDiagnostics(processing, lines, diagnostics); err != nil {
		t.Fatalf("validateDiagnostics(valid) error: %v", err)
	}

	tests := []struct {
		name   string
		change func(*[]Processing, *[]Line, *[]Diagnostic)
	}{
		{name: "пустой id", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].ID = "" }},
		{name: "повтор id", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[1].ID = (*d)[0].ID }},
		{name: "неизвестный source", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Source = "missing" }},
		{name: "пустой code", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Code = "" }},
		{name: "неизвестная severity", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Severity = Severity("notice") }},
		{name: "пустое message", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Message = "" }},
		{name: "неизвестный scope", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Scope = DiagnosticScope("file") }},
		{name: "fatal warning", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[1].Fatal = true }},
		{name: "element без location", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location = nil }},
		{name: "line на нескольких строках", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[1].Location.End.Line = 2 }},
		{name: "строка диапазона равна нулю", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location.Start.Line = 0 }},
		{name: "колонка равна нулю", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location.Start.Column = 0 }},
		{name: "обратный диапазон", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location.End.Column = 1 }},
		{name: "диапазон за концом строки", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location.End.Column = 99 }},
		{name: "related повторяет location", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) {
			(*d)[0].RelatedLocations = []Range{*(*d)[0].Location}
		}},
		{name: "повтор related", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) {
			value := Range{Start: Position{Line: 2, Column: 1}, End: Position{Line: 2, Column: 2}}
			(*d)[0].RelatedLocations = []Range{value, value}
		}},
		{name: "неизвестный error id", change: func(_ *[]Processing, l *[]Line, _ *[]Diagnostic) { (*l)[0].Elements[0].ErrorIDs = []string{"missing"} }},
		{name: "error id с line scope", change: func(_ *[]Processing, l *[]Line, _ *[]Diagnostic) { (*l)[0].Elements[0].ErrorIDs = []string{"d2"} }},
		{name: "нет обратной element ссылки", change: func(_ *[]Processing, l *[]Line, _ *[]Diagnostic) { (*l)[0].Elements[0].ErrorIDs = []string{} }},
		{name: "location не совпадает с элементом", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0].Location = rangePtr(1, 2, 1, 3) }},
		{name: "неверный line hasErrors", change: func(_ *[]Processing, l *[]Line, _ *[]Diagnostic) { (*l)[0].HasErrors = false }},
		{name: "warning устанавливает line hasErrors", change: func(_ *[]Processing, l *[]Line, d *[]Diagnostic) {
			(*d)[0].Severity = SeverityWarning
			(*l)[0].HasErrors = true
		}},
		{name: "нарушен порядок source", change: func(p *[]Processing, _ *[]Line, d *[]Diagnostic) {
			*p = append(*p, Processing{ID: "run-2", Tool: "validator"})
			(*d)[0].Source = "run-2"
		}},
		{name: "нарушен порядок location", change: func(_ *[]Processing, _ *[]Line, d *[]Diagnostic) { (*d)[0], (*d)[1] = (*d)[1], (*d)[0] }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotProcessing, gotLines, gotDiagnostics := validDiagnostics()
			test.change(&gotProcessing, &gotLines, &gotDiagnostics)
			assertInvariantError(t, validateDiagnostics(gotProcessing, gotLines, gotDiagnostics))
		})
	}
}

// TestValidateAcceptsCompleteResultAndDoesNotMutate проверяет общий путь и отсутствие нормализации модели.
func TestValidateAcceptsCompleteResultAndDoesNotMutate(t *testing.T) {
	result := completeResult()
	want := completeResult()
	if err := result.Validate(); err != nil {
		t.Fatalf("Validate() error: %v", err)
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("Validate() mutated result:\n got: %#v\nwant: %#v", result, want)
	}
}

// TestValidateStopsAtEveryInvalidSection проверяет вызов всех групп инвариантов общим методом.
func TestValidateStopsAtEveryInvalidSection(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Result)
	}{
		{name: "document", change: func(r *Result) { r.FormatVersion = "" }},
		{name: "processing", change: func(r *Result) { r.Processing[0].ID = "" }},
		{name: "lines", change: func(r *Result) { r.Lines[0].Number = 2 }},
		{name: "diagnostics", change: func(r *Result) { r.Diagnostics[0].Source = "missing" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := completeResult()
			test.change(&result)
			assertInvariantError(t, result.Validate())
		})
	}
}

// TestRuneSlice проверяет включительно-исключающие Unicode-колонки и пустые диапазоны.
func TestRuneSlice(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		start, end int
		want       string
		ok         bool
	}{
		{name: "ascii", raw: "abc", start: 1, end: 4, want: "abc", ok: true},
		{name: "кириллица", raw: "язык", start: 2, end: 4, want: "зы", ok: true},
		{name: "emoji одна колонка", raw: "a🙂b", start: 2, end: 3, want: "🙂", ok: true},
		{name: "составной символ две колонки", raw: "е\u0301", start: 1, end: 3, want: "е\u0301", ok: true},
		{name: "пусто в начале", raw: "abc", start: 1, end: 1, want: "", ok: true},
		{name: "пусто после строки", raw: "abc", start: 4, end: 4, want: "", ok: true},
		{name: "start ноль", raw: "abc", start: 0, end: 1, ok: false},
		{name: "обратный", raw: "abc", start: 3, end: 2, ok: false},
		{name: "за концом", raw: "abc", start: 1, end: 5, ok: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := runeSlice(test.raw, test.start, test.end)
			if got != test.want || ok != test.ok {
				t.Fatalf("runeSlice(%q, %d, %d) = %q, %v; want %q, %v", test.raw, test.start, test.end, got, ok, test.want, test.ok)
			}
		})
	}
}

// TestInvariantError проверяет стабильное человекочитаемое представление ошибки.
func TestInvariantError(t *testing.T) {
	err := &InvariantError{Path: "lines[1].elements[0].raw", Rule: "не совпадает с диапазоном"}
	message := err.Error()
	if !strings.Contains(message, err.Path) || !strings.Contains(message, err.Rule) {
		t.Fatalf("Error() = %q, want path and rule", message)
	}
}

// validResult возвращает минимальный допустимый результат с доступным источником.
func validResult() Result {
	return Result{
		FormatVersion: FormatVersion,
		Document: Document{
			DSLVersion: stringPtr("1.2"), FileName: stringPtr("lesson.txt"), FilePath: stringPtr("/course/lesson.txt"),
			Encoding: stringPtr("UTF-8"), HasBOM: boolPtr(false), LineCount: intPtr(0), ByteLength: intPtr(0),
			SHA256:   stringPtr("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			Metadata: Metadata{DocumentID: stringPtr("Lesson-1"), Title: stringPtr("Урок"), Subtitle: stringPtr("Введение"), Section: stringPtr("Основы"), Order: stringPtr("001"), ResourceDirs: []string{"assets"}},
		},
		Processing: []Processing{}, Lines: []Line{}, Diagnostics: []Diagnostic{},
	}
}

// validLines возвращает корректные связанные строки с Unicode-элементами.
func validLines() (Document, []Line) {
	lines := []Line{
		{Number: 1, Type: LineHeading, Raw: "# Урок🙂", EOL: EOLLF, Elements: []Element{
			{Type: ElementHeadingLevel, Raw: "#", Value: stringPtr("1"), Start: 1, End: 2, ErrorIDs: []string{}},
			{Type: ElementTitle, Raw: "Урок🙂", Value: stringPtr("Урок🙂"), Start: 3, End: 8, ErrorIDs: []string{}},
		}},
		{Number: 2, Type: LineContent, NestingLevel: 1, ParentLine: intPtr(1), Raw: "текст", EOL: EOLNone, Elements: []Element{
			{Type: ElementContent, Raw: "текст", Value: stringPtr("текст"), Start: 1, End: 6, ErrorIDs: []string{}},
		}},
	}
	return Document{LineCount: intPtr(len(lines))}, lines
}

// validDiagnostics возвращает согласованные диагностики element и line.
func validDiagnostics() ([]Processing, []Line, []Diagnostic) {
	processing := []Processing{{ID: "run-1", Tool: "dsl-parser"}}
	lines := []Line{
		{Number: 1, Type: LineInvalid, Raw: "?", EOL: EOLLF, HasErrors: true, Elements: []Element{{Type: ElementUnparsed, Raw: "?", Start: 1, End: 2, ErrorIDs: []string{"d1"}}}},
		{Number: 2, Type: LineContent, Raw: "x", EOL: EOLNone, Elements: []Element{{Type: ElementContent, Raw: "x", Value: stringPtr("x"), Start: 1, End: 2, ErrorIDs: []string{}}}},
	}
	diagnostics := []Diagnostic{
		{ID: "d1", Source: "run-1", Code: "P003", Severity: SeverityError, Message: "ошибка", Scope: ScopeElement, Location: rangePtr(1, 1, 1, 2), RelatedLocations: []Range{}},
		{ID: "d2", Source: "run-1", Code: "W001", Severity: SeverityWarning, Message: "предупреждение", Scope: ScopeLine, Location: rangePtr(2, 1, 2, 2), RelatedLocations: []Range{}},
	}
	return processing, lines, diagnostics
}

// completeResult возвращает корректный результат со связанными строками и диагностиками.
func completeResult() Result {
	processing, lines, diagnostics := validDiagnostics()
	return Result{
		FormatVersion: FormatVersion,
		Document:      Document{LineCount: intPtr(2), ByteLength: intPtr(3), Metadata: Metadata{ResourceDirs: []string{}}, HasErrors: true},
		Processing:    processing, Lines: lines, Diagnostics: diagnostics,
	}
}

// documentDiagnostic создаёт диагностику уровня документа.
func documentDiagnostic(severity Severity) Diagnostic {
	return Diagnostic{ID: "d", Source: "run", Code: "code", Severity: severity, Message: "message", Scope: ScopeDocument, RelatedLocations: []Range{}}
}

// assertInvariantError проверяет тип и содержательность ошибки инварианта.
func assertInvariantError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("validation error = nil")
	}
	var invariant *InvariantError
	if !errors.As(err, &invariant) {
		t.Fatalf("error type = %T, want *InvariantError", err)
	}
	if invariant.Path == "" || invariant.Rule == "" {
		t.Fatalf("InvariantError = %#v, want non-empty path and rule", invariant)
	}
}

// rangePtr создаёт диапазон с включительным началом и исключающим концом.
func rangePtr(startLine, startColumn, endLine, endColumn int) *Range {
	return &Range{Start: Position{Line: startLine, Column: startColumn}, End: Position{Line: endLine, Column: endColumn}}
}

// stringPtr возвращает указатель на строку.
func stringPtr(value string) *string { return &value }

// intPtr возвращает указатель на целое число.
func intPtr(value int) *int { return &value }

// boolPtr возвращает указатель на логическое значение.
func boolPtr(value bool) *bool { return &value }
