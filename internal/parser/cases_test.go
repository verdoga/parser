package parser

import (
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"
	"testing"

	"dslparser/internal/grammar"
	"dslparser/internal/grammar/v1_2"
	"dslparser/internal/model"
)

// TestDecodeSource проверяет точное выделение строк и первичные ошибки байтового слоя.
func TestDecodeSource(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		bom     bool
		lines   []physicalLine
		failure *sourceFailure
	}{
		{name: "empty", data: nil, lines: []physicalLine{}},
		{name: "one line without eol", data: []byte("text"), lines: []physicalLine{{number: 1, raw: "text", eol: model.EOLNone}}},
		{name: "final lf", data: []byte("one\n"), lines: []physicalLine{{number: 1, raw: "one", eol: model.EOLLF}}},
		{name: "mixed endings", data: []byte("one\r\ntwo\nthree"), lines: []physicalLine{{1, "one", model.EOLCRLF}, {2, "two", model.EOLLF}, {3, "three", model.EOLNone}}},
		{name: "leading bom", data: append([]byte{0xef, 0xbb, 0xbf}, []byte("текст\n")...), bom: true, lines: []physicalLine{{1, "текст", model.EOLLF}}},
		{name: "embedded bom is content", data: []byte("a\ufeffb"), lines: []physicalLine{{1, "a\ufeffb", model.EOLNone}}},
		{name: "invalid utf8", data: []byte{0xff}, lines: nil, failure: &sourceFailure{code: "P001", fatal: true}},
		{name: "lone cr", data: []byte("one\rtwo"), lines: []physicalLine{{1, "one", model.EOLNone}, {2, "two", model.EOLNone}}, failure: &sourceFailure{code: "P002", fatal: true, fragment: "\\r"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source, failure := decodeSource(test.data)
			if source.hasBOM != test.bom {
				t.Fatalf("hasBOM = %v, want %v", source.hasBOM, test.bom)
			}
			if !slices.Equal(source.lines, test.lines) {
				t.Fatalf("lines = %#v, want %#v", source.lines, test.lines)
			}
			if test.failure == nil {
				if failure != nil {
					t.Fatalf("failure = %#v, want nil", failure)
				}
				return
			}
			if failure == nil || failure.code != test.failure.code || failure.fatal != test.failure.fatal {
				t.Fatalf("failure = %#v, want code=%s fatal=%v", failure, test.failure.code, test.failure.fatal)
			}
		})
	}
}

// TestSplitLinesRejectsEveryLoneCR проверяет позицию каждого неподдерживаемого CR.
func TestSplitLinesRejectsEveryLoneCR(t *testing.T) {
	t.Parallel()

	for _, text := range []string{"\r", "a\rb", "a\n\rb", "a\r\nb\r"} {
		t.Run(fmt.Sprintf("%q", text), func(t *testing.T) {
			_, failure := splitLines(text)
			if failure == nil || failure.code != "P002" || failure.location == nil {
				t.Fatalf("failure = %#v, want located P002", failure)
			}
		})
	}
}

// TestUnicodeColumns проверяет перевод байтовых границ в колонки Unicode.
func TestUnicodeColumns(t *testing.T) {
	t.Parallel()

	raw := "Aйe\u0301😀"
	tests := []struct {
		byteOffset int
		column     int
	}{
		{0, 1}, {1, 2}, {3, 3}, {4, 4}, {6, 5}, {10, 6},
	}
	for _, test := range tests {
		if got := runeColumn(raw, test.byteOffset); got != test.column {
			t.Errorf("runeColumn(%q, %d) = %d, want %d", raw, test.byteOffset, got, test.column)
		}
	}

	start, end := runeRange(raw, 3, 10)
	if start != 3 || end != 6 {
		t.Fatalf("runeRange = (%d, %d), want (3, 6)", start, end)
	}
	start, end = runeRange(raw, 6, 6)
	if start != 5 || end != 5 {
		t.Fatalf("empty runeRange = (%d, %d), want (5, 5)", start, end)
	}
}

// TestTrimRecognitionSpace проверяет удаление только крайних SPACE и TAB.
func TestTrimRecognitionSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		raw        string
		text       string
		start, end int
	}{
		{"unchanged", "abc", "abc", 0, 3},
		{"ascii whitespace", " \tabc\t ", "abc", 2, 5},
		{"unicode is retained", "\u00a0abc\u00a0", "\u00a0abc\u00a0", 0, 7},
		{"only whitespace", " \t ", "", 3, 3},
		{"unicode columns use byte bounds", "  мир ", "мир", 2, 8},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := trimRecognitionSpace(test.raw)
			if got.text != test.text || got.byteStart != test.start || got.byteEnd != test.end {
				t.Fatalf("trimRecognitionSpace(%q) = %#v, want text=%q bounds=[%d,%d)", test.raw, got, test.text, test.start, test.end)
			}
		})
	}
}

// TestElementConstruction проверяет элементы, пустые диапазоны, ошибки и сортировку.
func TestElementConstruction(t *testing.T) {
	t.Parallel()

	value := "мир"
	line := physicalLine{number: 4, raw: "  @tag мир", eol: model.EOLNone}
	element := elementFromDecision(line, grammar.ElementDecision{Type: model.ElementContent, Value: &value, ByteStart: 7, ByteEnd: 13})
	if element.Type != model.ElementContent || element.Raw != "мир" || element.Value == nil || *element.Value != value || element.Start != 8 || element.End != 11 || element.ErrorIDs == nil {
		t.Fatalf("elementFromDecision = %#v", element)
	}
	value = "changed"
	if *element.Value != "мир" {
		t.Fatal("element value aliases grammar decision")
	}

	empty := emptyElement(line, 7, model.ElementIdentifier)
	if empty.Raw != "" || empty.Value != nil || empty.Start != 7 || empty.End != 7 || empty.ErrorIDs == nil {
		t.Fatalf("emptyElement = %#v", empty)
	}
	attachElementError(&empty, "d2")
	attachElementError(&empty, "d1")
	attachElementError(&empty, "d2")
	if !slices.Equal(empty.ErrorIDs, []string{"d2", "d1"}) {
		t.Fatalf("error IDs = %v", empty.ErrorIDs)
	}

	elements := []model.Element{{Start: 5, End: 6}, {Start: 1, End: 3}, {Start: 1, End: 2}}
	sortElements(elements)
	if elements[0].End != 2 || elements[1].End != 3 || elements[2].Start != 5 {
		t.Fatalf("sortElements = %#v", elements)
	}
}

// TestParserState проверяет независимые области состояния, стек и выбор родителя.
func TestParserState(t *testing.T) {
	t.Parallel()

	state := newParseState()
	if state.blocks == nil {
		t.Fatal("newParseState returned a nil block stack")
	}
	if _, ok := state.topBlock(); ok {
		t.Fatal("empty state unexpectedly has a top block")
	}
	frame := blockFrame{tag: "text", openingLine: 2, mode: grammar.ContentOpaque, context: "text"}
	state.openBlock(frame)
	if got, ok := state.topBlock(); !ok || got != frame {
		t.Fatalf("topBlock = %#v, %v", got, ok)
	}
	if got := state.parentFor(grammar.Parents(grammar.ParentBlock)); got != 2 {
		t.Fatalf("block parent = %d, want 2", got)
	}
	if got := state.closeBlock(); got != frame {
		t.Fatalf("closeBlock = %#v, want %#v", got, frame)
	}

	state.applyTransition(grammar.Transition{SetStep: true}, 3)
	state.applyTransition(grammar.Transition{OpenTask: true}, 4)
	state.applyTransition(grammar.Transition{SetInnerStep: true}, 5)
	state.applyTransition(grammar.Transition{OpenVariants: true}, 6)
	state.applyTransition(grammar.Transition{SetVariant: true}, 7)
	parent := grammar.Parents(grammar.ParentVariant, grammar.ParentInnerStep, grammar.ParentTask, grammar.ParentStep)
	if got := state.parentFor(parent); got != 7 {
		t.Fatalf("preferred parent = %d, want 7", got)
	}
	state.applyTransition(grammar.Transition{ClearVariants: true}, 8)
	if state.variants != 0 || state.variant != 0 || state.innerStep != 5 {
		t.Fatalf("ClearVariants changed wrong state: %#v", state)
	}
	state.applyTransition(grammar.Transition{CloseTask: true}, 9)
	if state.activeTask != 0 || state.innerStep != 0 || state.variants != 0 || state.variant != 0 || state.currentStep != 3 {
		t.Fatalf("CloseTask changed wrong state: %#v", state)
	}
}

// TestLexicalTail проверяет безопасное представление необработанного хвоста.
func TestLexicalTail(t *testing.T) {
	t.Parallel()

	source := []physicalLine{{1, "ignored", model.EOLLF}, {2, " \t ", model.EOLCRLF}, {3, `\@raw {`, model.EOLNone}}
	got := appendLexicalTail([]model.Line{{Number: 1}}, source, 1)
	if len(got) != 3 || got[1].Type != model.LineBlank || got[2].Type != model.LineContent {
		t.Fatalf("appendLexicalTail = %#v", got)
	}
	if got[2].Raw != `\@raw {` || got[2].ParentLine != nil || got[2].NestingLevel != 0 || len(got[2].Elements) != 1 {
		t.Fatalf("tail content line = %#v", got[2])
	}
	element := got[2].Elements[0]
	if element.Type != model.ElementContent || element.Raw != `\@raw {` || element.Value == nil || *element.Value != `\@raw {` || element.Start != 1 || element.End != 8 {
		t.Fatalf("tail element = %#v", element)
	}
}

// TestMetadata проверяет однозначность метаданных, порядок путей и владение памятью.
func TestMetadata(t *testing.T) {
	t.Parallel()

	collector := metadataCollector{}
	collector.documentIDs = []locatedValue{{"Doc-A", 2}}
	collector.titles = []locatedValue{{"Title", 3}}
	collector.subtitles = []locatedValue{{"Subtitle", 4}}
	collector.sections = []locatedValue{{"01", 5}}
	collector.orders = []locatedValue{{"007", 6}}
	collector.resourceDirs = []locatedValue{{"first", 7}, {"second", 7}}

	metadata := collector.metadata(true)
	if metadata.DocumentID == nil || *metadata.DocumentID != "Doc-A" || metadata.Order == nil || *metadata.Order != "007" {
		t.Fatalf("metadata = %#v", metadata)
	}
	if !slices.Equal(metadata.ResourceDirs, []string{"first", "second"}) {
		t.Fatalf("resource dirs = %v", metadata.ResourceDirs)
	}
	collector.resourceDirs[0].value = "changed"
	if metadata.ResourceDirs[0] != "first" {
		t.Fatal("metadata aliases collector paths")
	}
	if unavailable := collector.metadata(false); unavailable.ResourceDirs != nil {
		t.Fatalf("unavailable resource dirs = %#v, want nil", unavailable.ResourceDirs)
	}
	collector.documentIDs = append(collector.documentIDs, locatedValue{"Doc-B", 8})
	if got := collector.metadata(true).DocumentID; got != nil {
		t.Fatalf("ambiguous document ID = %q, want nil", *got)
	}
	if got := uniqueValue(nil); got != nil {
		t.Fatalf("uniqueValue(nil) = %q", *got)
	}
}

// TestMetadataObserve проверяет сбор метаданных только из подходящих типизированных элементов.
func TestMetadataObserve(t *testing.T) {
	t.Parallel()

	collector := metadataCollector{}
	lines := []model.Line{
		{Number: 2, Type: model.LineTag, Elements: []model.Element{testElement(model.ElementTag, "document-id"), testElement(model.ElementIdentifier, "Doc-A")}},
		{Number: 3, Type: model.LineHeading, Elements: []model.Element{testElement(model.ElementHeadingLevel, "1"), testElement(model.ElementTitle, "Title")}},
		{Number: 4, Type: model.LineHeading, Elements: []model.Element{testElement(model.ElementHeadingLevel, "2"), testElement(model.ElementTitle, "Subtitle")}},
		{Number: 5, Type: model.LineTag, Elements: []model.Element{testElement(model.ElementTag, "section"), testElement(model.ElementNumber, "01")}},
		{Number: 6, Type: model.LineTag, Elements: []model.Element{testElement(model.ElementTag, "order"), testElement(model.ElementNumber, "007")}},
		{Number: 7, Type: model.LineTag, Elements: []model.Element{testElement(model.ElementTag, "resource-dir"), testElement(model.ElementResourcePath, "one"), testElement(model.ElementResourcePath, "two")}},
		{Number: 8, Type: model.LineContent, Elements: []model.Element{testElement(model.ElementIdentifier, "ignored")}},
	}
	for _, line := range lines {
		collector.observe(line)
	}
	got := collector.metadata(true)
	if got.DocumentID == nil || *got.DocumentID != "Doc-A" || got.Title == nil || *got.Title != "Title" || got.Subtitle == nil || *got.Subtitle != "Subtitle" || got.Section == nil || *got.Section != "01" || got.Order == nil || *got.Order != "007" || !slices.Equal(got.ResourceDirs, []string{"one", "two"}) {
		t.Fatalf("observed metadata = %#v", got)
	}
}

// TestRecoveryDecision проверяет продолжение только при явном однозначном решении.
func TestRecoveryDecision(t *testing.T) {
	t.Parallel()
	if !canRecover(grammar.RecoveryDecision{Continue: true}) {
		t.Fatal("Continue=true was rejected")
	}
	if canRecover(grammar.RecoveryDecision{Continue: false}) {
		t.Fatal("Continue=false was accepted")
	}
}

// TestBuildDocument проверяет доступные сведения документа и независимые указатели.
func TestBuildDocument(t *testing.T) {
	t.Parallel()

	version := "1.2"
	metadata := model.Metadata{DocumentID: stringPointer("Doc"), ResourceDirs: []string{"assets"}}
	document := buildDocument(parserInput([]byte("text")), physicalSource{hasBOM: true, lines: []physicalLine{{number: 1, raw: "text"}}}, &version, metadata)
	if document.DSLVersion == nil || *document.DSLVersion != "1.2" || document.FileName == nil || *document.FileName != "lesson.txt" || document.FilePath == nil || *document.FilePath != "/course/lesson.txt" || document.Encoding == nil || *document.Encoding != "UTF-8" || document.HasBOM == nil || !*document.HasBOM || document.LineCount == nil || *document.LineCount != 1 || document.ByteLength == nil || *document.ByteLength != 4 || document.SHA256 == nil {
		t.Fatalf("buildDocument = %#v", document)
	}
	version = "changed"
	*metadata.DocumentID = "changed"
	metadata.ResourceDirs[0] = "changed"
	if *document.DSLVersion != "1.2" || *document.Metadata.DocumentID != "Doc" || document.Metadata.ResourceDirs[0] != "assets" {
		t.Fatal("document aliases its inputs")
	}
}

// TestProbeVersion проверяет строгую форму первой строки и Unicode-диапазоны.
func TestProbeVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, raw, value, problem string
	}{
		{"valid", "@dsl-version 1.2", "1.2", ""},
		{"trimmed", " \t@DSL-VERSION 9.9\t ", "9.9", ""},
		{"tab separator", "@dsl-version\t1.2", "", "P013"},
		{"missing value", "@dsl-version ", "", "P013"},
		{"extra value", "@dsl-version 1.2 extra", "", "P013"},
		{"other tag", "@document-id x", "", "P013"},
		{"blank", "", "", "P013"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := probeVersion(physicalLine{number: 1, raw: test.raw})
			if test.problem != "" {
				if got.value != nil || got.problem == nil || string(got.problem.Code) != test.problem {
					t.Fatalf("probeVersion(%q) = %#v", test.raw, got)
				}
				return
			}
			if got.problem != nil || got.value == nil || *got.value != test.value || got.element == nil || got.element.Type != model.ElementVersion {
				t.Fatalf("probeVersion(%q) = %#v", test.raw, got)
			}
		})
	}
}

// TestProbeSource проверяет лёгкое извлечение версии и единственного document-id.
func TestProbeSource(t *testing.T) {
	t.Parallel()

	registry := parserRegistry(t)
	tests := []struct {
		name, text, version, documentID string
	}{
		{"valid", "@dsl-version 1.2\n@document-id Doc-A\n", "1.2", "Doc-A"},
		{"unsupported version is retained", "@dsl-version 9.9\n@document-id Doc-A", "9.9", ""},
		{"missing version", "@document-id Doc-A", "", ""},
		{"duplicate id", "@dsl-version 1.2\n@document-id A\n@document-id B", "1.2", ""},
		{"missing id", "@dsl-version 1.2\ntext", "1.2", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ProbeSource([]byte(test.text), registry)
			if err != nil {
				t.Fatalf("ProbeSource: %v", err)
			}
			assertOptionalString(t, "DSLVersion", got.DSLVersion, test.version)
			assertOptionalString(t, "DocumentID", got.DocumentID, test.documentID)
		})
	}
}

// TestParseSourceFailures проверяет фатальные исходы до структурного разбора.
func TestParseSourceFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data []byte
		code string
	}{
		{"invalid utf8", []byte{0xff}, "P001"},
		{"lone cr", []byte("@dsl-version 1.2\rtext"), "P002"},
		{"empty", nil, "P013"},
		{"missing version", []byte("text"), "P013"},
		{"unsupported version", []byte("@dsl-version 9.9\ntext"), "P014"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Parse(parserInput(test.data), parserRegistry(t))
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if !result.Fatal || len(result.Diagnostics) != 1 || result.Diagnostics[0].Code != test.code || !result.Diagnostics[0].Fatal {
				t.Fatalf("result diagnostics = %#v, fatal=%v", result.Diagnostics, result.Fatal)
			}
			if test.code == "P001" && (result.Document.Encoding != nil || result.Document.LineCount != nil || len(result.Lines) != 0) {
				t.Fatalf("invalid UTF-8 exposed decoded data: %#v", result)
			}
			if test.code == "P002" && len(result.Lines) == 0 {
				t.Fatal("P002 did not preserve lexical lines")
			}
		})
	}
}

// TestParseEveryTagForm проверяет репрезентативные формы всех 34 тегов DSL v1.2.
func TestParseEveryTagForm(t *testing.T) {
	t.Parallel()

	lines := []string{
		"@dsl-version 1.2", "@document-id Doc", "@section 01", "@order 007", `@resource-dir assets, "two words"`,
		"# Title", "## Subtitle", "### Step", "@header Header", "@task task-1", "@step Inner", "@speaking", "@newpage", "@endtask",
		"@editor Editor text", `@media audio "clip.mp3"`, "@example Example", "@wordlist one; two", "@script Script", "@key Key", "@note Note", "@alt Alt", "@hint Hint", "@include file.txt", "@answer Answer", "@question Question",
		"@text {", "text", "}", "@table {", "row", "}", "@instr {", "instruction", "}", "@fragment {", "fragment", "}", "@multifill {", "fill", "}", "@choice {", "choice", "}", "@multichoice {", "choice", "}", "@matching {", "match", "}", "@ordering {", "order", "}", "@variants", "@variant A",
	}
	result := parseText(t, strings.Join(lines, "\n"))
	if result.Fatal || len(result.Diagnostics) != 0 {
		t.Fatalf("valid forms produced diagnostics: %#v", result.Diagnostics)
	}
	if len(result.Lines) != len(lines) {
		t.Fatalf("line count = %d, want %d", len(result.Lines), len(lines))
	}
	for i, line := range result.Lines {
		if line.Number != i+1 || line.Raw != lines[i] || line.Elements == nil {
			t.Errorf("line %d = %#v", i+1, line)
		}
	}
	if result.Document.Metadata.DocumentID == nil || *result.Document.Metadata.DocumentID != "Doc" || result.Document.Metadata.Order == nil || *result.Document.Metadata.Order != "007" {
		t.Fatalf("metadata = %#v", result.Document.Metadata)
	}
}

// TestOpaqueContent проверяет непрозрачные режимы и отсутствие placeholder-элементов.
func TestOpaqueContent(t *testing.T) {
	t.Parallel()

	text := "@dsl-version 1.2\n@text {\n<b>@task x</b>\n\n}\n@table {\n---\n}\n@multifill {\n\n_____{answer}\n}\n"
	result := parseText(t, text)
	if result.Fatal {
		t.Fatalf("Parse is fatal: %#v", result.Diagnostics)
	}
	wants := map[int]model.LineType{3: model.LineContent, 4: model.LineBlank, 7: model.LineContent, 10: model.LineContent, 11: model.LineContent}
	for number, want := range wants {
		if got := result.Lines[number-1].Type; got != want {
			t.Errorf("line %d type = %q, want %q", number, got, want)
		}
	}
	for _, line := range result.Lines {
		for _, element := range line.Elements {
			if element.Type == model.ElementPlaceholder {
				t.Errorf("line %d contains forbidden placeholder", line.Number)
			}
		}
	}
}

// TestParserRecovery проверяет восстановимые ошибки, фатальный хвост и незакрытый блок.
func TestParserRecovery(t *testing.T) {
	t.Parallel()

	recoverable := []struct {
		name, line, code string
	}{
		{"unknown tag", "@unknown", "P003"},
		{"missing separator", "@taskid", "P004"},
		{"unsupported form", "@speaking {", "P005"},
		{"missing argument", "@task", "P006"},
		{"extra content", "@newpage extra", "P007"},
		{"unexpected close", "}", "P009"},
		{"malformed close", "} trailing", "P010"},
		{"unescaped brace", "text { brace", "P012"},
	}
	for _, test := range recoverable {
		t.Run(test.name, func(t *testing.T) {
			result := parseText(t, "@dsl-version 1.2\n"+test.line+"\nafter")
			assertDiagnostic(t, result, test.code, false)
			if result.Fatal || len(result.Lines) != 3 {
				t.Fatalf("recovery result = %#v", result)
			}
		})
	}

	t.Run("ambiguous block opening leaves lexical tail", func(t *testing.T) {
		result := parseText(t, "@dsl-version 1.2\n@text { trailing\n\\@literal\n\n")
		assertDiagnostic(t, result, "P008", true)
		if result.Lines[2].Type != model.LineContent || result.Lines[2].Elements[0].Raw != `\@literal` || result.Lines[3].Type != model.LineBlank {
			t.Fatalf("lexical tail = %#v", result.Lines[2:])
		}
	})

	t.Run("unclosed block", func(t *testing.T) {
		result := parseText(t, "@dsl-version 1.2\n@text {\ncontent")
		assertDiagnostic(t, result, "P011", true)
		if len(result.Diagnostics[0].RelatedLocations) != 1 {
			t.Fatalf("related locations = %#v", result.Diagnostics[0].RelatedLocations)
		}
	})
}

// TestParserDiagnosticIntegration проверяет реестр, дедупликацию и обратные ссылки элементов.
func TestParserDiagnosticIntegration(t *testing.T) {
	t.Parallel()

	result := parseText(t, "@dsl-version 1.2\n@task\nafter")
	assertDiagnostic(t, result, "P006", false)
	diagnostic := result.Diagnostics[0]
	if diagnostic.Source != "run-parser-tests" || diagnostic.ID == "" || diagnostic.Location == nil || diagnostic.Scope != model.ScopeElement {
		t.Fatalf("diagnostic = %#v", diagnostic)
	}
	count := 0
	for _, element := range result.Lines[1].Elements {
		for _, id := range element.ErrorIDs {
			if id == diagnostic.ID {
				count++
			}
		}
	}
	if count != 1 {
		t.Fatalf("diagnostic %q has %d element back-references, want 1", diagnostic.ID, count)
	}
}

// TestParseRejectsMismatchedSHA256 проверяет контракт точного отпечатка входных байтов.
func TestParseRejectsMismatchedSHA256(t *testing.T) {
	t.Parallel()

	input := parserInput([]byte("@dsl-version 1.2"))
	input.SHA256 = strings.Repeat("0", 64)
	if _, err := Parse(input, parserRegistry(t)); err == nil {
		t.Fatal("Parse accepted a mismatched SHA-256")
	}
}

// TestParserParentAndNesting проверяет логических родителей независимо от отступов.
func TestParserParentAndNesting(t *testing.T) {
	t.Parallel()

	result := parseText(t, "@dsl-version 1.2\n### Outer\n  @task t\n\t@step Inner\ncontent\n@endtask\ncontent")
	if result.Fatal || len(result.Diagnostics) != 0 {
		t.Fatalf("Parse diagnostics = %#v", result.Diagnostics)
	}
	wants := []struct{ line, parent, nesting int }{{2, 0, 0}, {3, 2, 1}, {4, 3, 2}, {5, 4, 3}, {6, 3, 2}, {7, 2, 1}}
	for _, want := range wants {
		line := result.Lines[want.line-1]
		parent := 0
		if line.ParentLine != nil {
			parent = *line.ParentLine
		}
		if parent != want.parent || line.NestingLevel != want.nesting {
			t.Errorf("line %d parent/nesting = %d/%d, want %d/%d", want.line, parent, line.NestingLevel, want.parent, want.nesting)
		}
	}
}

// TestParseDoesNotAliasInput проверяет владение входными байтами и полями результата.
func TestParseDoesNotAliasInput(t *testing.T) {
	t.Parallel()

	data := []byte("@dsl-version 1.2\n@document-id Doc")
	result, err := Parse(parserInput(data), parserRegistry(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	for i := range data {
		data[i] = 'x'
	}
	if result.Lines[0].Raw != "@dsl-version 1.2" || result.Document.Metadata.DocumentID == nil || *result.Document.Metadata.DocumentID != "Doc" {
		t.Fatalf("result aliases input: %#v", result)
	}
}

// parserRegistry создаёт реестр нормативной версии или завершает тест.
func parserRegistry(t *testing.T) grammar.Registry {
	t.Helper()
	registry, err := grammar.NewRegistry(v1_2.New())
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	return registry
}

// parserInput создаёт согласованный вход парсера.
func parserInput(data []byte) Input {
	sum := sha256.Sum256(data)
	return Input{Bytes: data, FileName: "lesson.txt", FilePath: "/course/lesson.txt", SHA256: fmt.Sprintf("%x", sum), ProcessingID: "run-parser-tests"}
}

// parseText разбирает текст нормативной грамматикой или завершает тест при контрактной ошибке.
func parseText(t *testing.T, text string) Result {
	t.Helper()
	result, err := Parse(parserInput([]byte(text)), parserRegistry(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return result
}

// assertOptionalString сравнивает необязательную строку с пустым обозначением nil.
func assertOptionalString(t *testing.T, name string, got *string, want string) {
	t.Helper()
	if want == "" {
		if got != nil {
			t.Fatalf("%s = %q, want nil", name, *got)
		}
		return
	}
	if got == nil || *got != want {
		t.Fatalf("%s = %v, want %q", name, got, want)
	}
}

// assertDiagnostic проверяет единственную ожидаемую диагностику и итоговую фатальность.
func assertDiagnostic(t *testing.T, result Result, code string, fatal bool) {
	t.Helper()
	if len(result.Diagnostics) != 1 || result.Diagnostics[0].Code != code || result.Diagnostics[0].Fatal != fatal || result.Fatal != fatal {
		t.Fatalf("diagnostics = %#v, result.Fatal = %v; want %s fatal=%v", result.Diagnostics, result.Fatal, code, fatal)
	}
}

// testElement создаёт минимальный элемент для проверки сборщика метаданных.
func testElement(elementType model.ElementType, value string) model.Element {
	return model.Element{Type: elementType, Value: stringPointer(value), ErrorIDs: []string{}}
}

// stringPointer возвращает указатель на независимую строку.
func stringPointer(value string) *string {
	return &value
}
