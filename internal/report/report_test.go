package report

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"dslparser/internal/model"
)

// TestBuildOwnershipAndArrays проверяет глубокое копирование входа и представление обязательных массивов.
func TestBuildOwnershipAndArrays(t *testing.T) {
	dslVersion := "1.2"
	fileName := "lesson.txt"
	filePath := "/course/lesson.txt"
	encoding := "UTF-8"
	hasBOM := false
	lineCount := 1
	byteLength := 10
	sha256 := strings.Repeat("a", 64)
	documentID := "Lesson"
	title := "Заголовок"
	subtitle := "Подзаголовок"
	section := "Раздел"
	order := "01"
	version := "test"
	startedAt := "2026-09-17T12:00:00Z"
	duration := 4
	parent := 1
	value := "текст"
	location := model.Range{Start: model.Position{Line: 2, Column: 1}, End: model.Position{Line: 2, Column: 6}}

	input := Input{
		Document: model.Document{
			DSLVersion: &dslVersion, FileName: &fileName, FilePath: &filePath,
			Encoding: &encoding, HasBOM: &hasBOM, LineCount: &lineCount,
			ByteLength: &byteLength, SHA256: &sha256,
			Metadata: model.Metadata{
				DocumentID: &documentID, Title: &title, Subtitle: &subtitle,
				Section: &section, Order: &order, ResourceDirs: []string{"assets"},
			},
		},
		Processing: model.Processing{
			ID: "run-1", Tool: "dsl-parser", Version: &version,
			StartedAt: &startedAt, DurationMS: &duration,
		},
		Lines: []model.Line{{
			Number: 2, Type: model.LineContent, NestingLevel: 1, ParentLine: &parent,
			Raw: "текст", EOL: model.EOLNone,
			Elements: []model.Element{{
				Type: model.ElementContent, Raw: "текст", Value: &value,
				Start: 1, End: 6, ErrorIDs: []string{"diagnostic-1"},
			}},
		}},
		Diagnostics: []model.Diagnostic{{
			ID: "diagnostic-1", Source: "run-1", Code: "P003",
			Severity: model.SeverityWarning, Message: "сообщение", Scope: model.ScopeElement,
			Location: &location, RelatedLocations: []model.Range{location},
		}},
	}

	result := Build(input)
	if result.FormatVersion != model.FormatVersion {
		t.Fatalf("FormatVersion = %q, ожидалось %q", result.FormatVersion, model.FormatVersion)
	}
	if len(result.Processing) != 1 || result.Processing[0].ID != input.Processing.ID {
		t.Fatalf("Processing = %#v, ожидалась одна входная попытка", result.Processing)
	}

	// Изменение каждого вида изменяемых входных данных не должно проникать в результат.
	dslVersion = "changed"
	documentID = "changed"
	input.Document.Metadata.ResourceDirs[0] = "changed"
	version = "changed"
	parent = 99
	value = "changed"
	input.Lines[0].Elements[0].ErrorIDs[0] = "changed"
	input.Diagnostics[0].Location.Start.Line = 99
	input.Diagnostics[0].RelatedLocations[0].Start.Line = 99
	input.Lines[0].Raw = "changed"
	input.Diagnostics[0].Message = "changed"

	if *result.Document.DSLVersion != "1.2" || *result.Document.Metadata.DocumentID != "Lesson" {
		t.Errorf("указатели документа не скопированы: %#v", result.Document)
	}
	if got := result.Document.Metadata.ResourceDirs[0]; got != "assets" {
		t.Errorf("ResourceDirs[0] = %q, ожидалось assets", got)
	}
	if *result.Processing[0].Version != "test" {
		t.Errorf("Processing.Version = %q, ожидалось test", *result.Processing[0].Version)
	}
	if *result.Lines[0].ParentLine != 1 || *result.Lines[0].Elements[0].Value != "текст" {
		t.Errorf("строка не скопирована глубоко: %#v", result.Lines[0])
	}
	if got := result.Lines[0].Elements[0].ErrorIDs[0]; got != "diagnostic-1" {
		t.Errorf("ErrorIDs[0] = %q, ожидалось diagnostic-1", got)
	}
	if result.Diagnostics[0].Location.Start.Line != 2 || result.Diagnostics[0].RelatedLocations[0].Start.Line != 2 {
		t.Errorf("диапазоны диагностики не скопированы: %#v", result.Diagnostics[0])
	}
	if result.Lines[0].Raw != "текст" || result.Diagnostics[0].Message != "сообщение" {
		t.Error("изменение входных срезов изменило результат")
	}

	result = Build(Input{})
	if result.Document.Metadata.ResourceDirs == nil || result.Processing == nil || result.Lines == nil || result.Diagnostics == nil {
		t.Fatalf("обязательные массивы должны быть пустыми, не nil: %#v", result)
	}
}

// TestBuildHasErrors проверяет правила вычисления признаков ошибок документа и строк.
func TestBuildHasErrors(t *testing.T) {
	rangeAt := func(startLine, endLine int) *model.Range {
		return &model.Range{
			Start: model.Position{Line: startLine, Column: 1},
			End:   model.Position{Line: endLine, Column: 2},
		}
	}
	tests := []struct {
		name         string
		diagnostics  []model.Diagnostic
		errorIDs     []string
		wantLine     bool
		wantDocument bool
	}{
		{name: "нет диагностик"},
		{name: "warning в строке", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityWarning, Location: rangeAt(2, 2)}}},
		{name: "recommendation в строке", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityRecommendation, Location: rangeAt(2, 2)}}},
		{name: "error начинается в строке", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError, Location: rangeAt(2, 2)}}, wantLine: true, wantDocument: true},
		{name: "многострочный error начинается в строке", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError, Location: rangeAt(2, 3)}}, wantLine: true, wantDocument: true},
		{name: "error дочерней строки не наследуется", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError, Location: rangeAt(3, 3)}}, wantDocument: true},
		{name: "document error без location", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError}}, wantDocument: true},
		{name: "element error по ссылке", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError, Location: rangeAt(3, 3)}}, errorIDs: []string{"d1"}, wantLine: true, wantDocument: true},
		{name: "ссылка на warning", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityWarning, Location: rangeAt(3, 3)}}, errorIDs: []string{"d1"}},
		{name: "неизвестная ссылка", diagnostics: []model.Diagnostic{{ID: "d1", Severity: model.SeverityError, Location: rangeAt(3, 3)}}, errorIDs: []string{"missing"}, wantDocument: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Build(Input{
				Lines: []model.Line{{
					Number:   2,
					Elements: []model.Element{{ErrorIDs: test.errorIDs}},
				}},
				Diagnostics: test.diagnostics,
			})
			if result.Document.HasErrors != test.wantDocument {
				t.Errorf("Document.HasErrors = %t, ожидалось %t", result.Document.HasErrors, test.wantDocument)
			}
			if result.Lines[0].HasErrors != test.wantLine {
				t.Errorf("Line.HasErrors = %t, ожидалось %t", result.Lines[0].HasErrors, test.wantLine)
			}
		})
	}
}

// TestCloneOptionalValues проверяет nil и независимость копий необязательных скалярных значений и диапазона.
func TestCloneOptionalValues(t *testing.T) {
	if cloneString(nil) != nil || cloneInt(nil) != nil || cloneBool(nil) != nil || cloneRange(nil) != nil {
		t.Fatal("копирование nil должно возвращать nil")
	}
	stringValue, intValue, boolValue := "value", 7, true
	rangeValue := model.Range{Start: model.Position{Line: 1, Column: 2}, End: model.Position{Line: 1, Column: 3}}
	stringCopy := cloneString(&stringValue)
	intCopy := cloneInt(&intValue)
	boolCopy := cloneBool(&boolValue)
	rangeCopy := cloneRange(&rangeValue)
	if stringCopy == &stringValue || intCopy == &intValue || boolCopy == &boolValue || rangeCopy == &rangeValue {
		t.Fatal("копии не должны разделять адреса с исходными значениями")
	}
	stringValue, intValue, boolValue = "changed", 8, false
	rangeValue.Start.Line = 9
	if *stringCopy != "value" || *intCopy != 7 || !*boolCopy || rangeCopy.Start.Line != 1 {
		t.Fatal("изменение исходных значений повлияло на копии")
	}
}

// TestMarshal проверяет валидацию, форматирование JSON, экранирование и порядок массивов.
func TestMarshal(t *testing.T) {
	t.Run("корректный результат", func(t *testing.T) {
		result := model.Result{
			FormatVersion: model.FormatVersion,
			Document:      model.Document{Metadata: model.Metadata{ResourceDirs: []string{}}},
			Processing:    []model.Processing{},
			Lines:         []model.Line{},
			Diagnostics:   []model.Diagnostic{},
		}
		data, err := Marshal(result)
		if err != nil {
			t.Fatalf("Marshal() вернул ошибку: %v", err)
		}
		if !bytes.HasSuffix(data, []byte("\n")) || bytes.HasSuffix(data, []byte("\n\n")) {
			t.Fatalf("результат должен завершаться ровно одним LF: %q", data)
		}
		if !bytes.Contains(data, []byte("\n  \"document\":")) || bytes.Contains(data, []byte("\\u003c")) {
			t.Fatalf("неожиданное форматирование JSON: %s", data)
		}
		var decoded model.Result
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("результат не является JSON: %v", err)
		}
	})

	t.Run("символы HTML и нормативный порядок", func(t *testing.T) {
		html := "<tag>&value>"
		result := model.Result{
			FormatVersion: model.FormatVersion,
			Document:      model.Document{Metadata: model.Metadata{ResourceDirs: []string{}}},
			Processing: []model.Processing{
				{ID: "first", Tool: html},
				{ID: "second", Tool: "tool-2"},
			},
			Lines: []model.Line{}, Diagnostics: []model.Diagnostic{},
		}
		data, err := Marshal(result)
		if err != nil {
			t.Fatalf("Marshal() вернул ошибку: %v", err)
		}
		if !bytes.Contains(data, []byte(html)) || bytes.Contains(data, []byte("\\u003c")) || bytes.Contains(data, []byte("\\u0026")) {
			t.Fatalf("HTML-символы экранированы: %s", data)
		}
		if bytes.Index(data, []byte(`"id": "first"`)) >= bytes.Index(data, []byte(`"id": "second"`)) {
			t.Fatalf("порядок processing изменён: %s", data)
		}
	})

	t.Run("нарушение инварианта", func(t *testing.T) {
		_, err := Marshal(model.Result{})
		if err == nil {
			t.Fatal("Marshal() не вернул ошибку валидации")
		}
		var invariantError *model.InvariantError
		if !errors.As(err, &invariantError) {
			t.Fatalf("ошибка %T не сохраняет причину InvariantError", err)
		}
	})
}
