package diagnostics

import (
	"errors"
	"reflect"
	"testing"

	"dslparser/internal/model"
)

// TestRegistry проверяет порядок, метаданные, ограничения и сообщения всех записей реестра.
func TestRegistry(t *testing.T) {
	readFailure := errors.New("permission denied")
	tests := []struct {
		registryCase
		title         string
		allowedScopes []model.DiagnosticScope
		details       Details
		message       string
	}{
		{registryCase: registryCase{name: "P001", code: P001, scope: model.ScopeDocument, fatalPolicy: fatalAlways}, title: "Некорректная кодировка UTF-8", allowedScopes: []model.DiagnosticScope{model.ScopeDocument}, message: "Некорректная кодировка UTF-8"},
		{registryCase: registryCase{name: "P002", code: P002, scope: model.ScopeLine, fatalPolicy: fatalAlways}, title: "Неподдерживаемый перевод строки", allowedScopes: []model.DiagnosticScope{model.ScopeLine}, message: "Неподдерживаемый перевод строки"},
		{registryCase: registryCase{name: "P003", code: P003, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Неизвестная DSL-конструкция", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Fragment: "@unknown"}, message: "Неизвестная DSL-конструкция: @unknown"},
		{registryCase: registryCase{name: "P004", code: P004, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Нарушен обязательный разделитель объявления", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Required: "пробел U+0020"}, message: "Нарушен обязательный разделитель объявления: требуется пробел U+0020"},
		{registryCase: registryCase{name: "P005", code: P005, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Использована неподдерживаемая форма тега", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Tag: "question", Received: "блочная"}, message: "Использована неподдерживаемая форма тега @question: блочная"},
		{registryCase: registryCase{name: "P006", code: P006, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Отсутствует обязательный структурный параметр", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Tag: "task", Required: "ID"}, message: "У тега @task отсутствует обязательный структурный параметр: ID"},
		{registryCase: registryCase{name: "P007", code: P007, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Лишние параметры или запрещённое содержимое объявления", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Fragment: "extra"}, message: "Лишние параметры или запрещённое содержимое объявления: extra"},
		{registryCase: registryCase{name: "P008", code: P008, scope: model.ScopeLine, fatalPolicy: fatalOnAmbiguity}, title: "Неправильное открытие блока", allowedScopes: []model.DiagnosticScope{model.ScopeLine}, details: Details{Fragment: "@text { content"}, message: "Неправильное открытие блока: @text { content"},
		{registryCase: registryCase{name: "P009", code: P009, scope: model.ScopeLine, fatalPolicy: fatalOnAmbiguity}, title: "Закрытие несуществующего блока", allowedScopes: []model.DiagnosticScope{model.ScopeLine}, message: "Закрытие несуществующего блока"},
		{registryCase: registryCase{name: "P010", code: P010, scope: model.ScopeLine, fatalPolicy: fatalOnAmbiguity}, title: "Неправильная строка закрытия блока", allowedScopes: []model.DiagnosticScope{model.ScopeLine}, details: Details{Fragment: "} text"}, message: "Неправильная строка закрытия блока: } text"},
		{registryCase: registryCase{name: "P011", code: P011, scope: model.ScopeBlock, fatalPolicy: fatalAlways}, title: "Незакрытый блок", allowedScopes: []model.DiagnosticScope{model.ScopeBlock}, details: Details{Tag: "text"}, message: "Незакрытый блок @text"},
		{registryCase: registryCase{name: "P012", code: P012, scope: model.ScopeElement, fatalPolicy: fatalOnAmbiguity}, title: "Неэкранированная фигурная скобка в текстовом значении", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Fragment: "{"}, message: "Неэкранированная фигурная скобка в текстовом значении: {"},
		{registryCase: registryCase{name: "P013", code: P013, scope: model.ScopeLine, fatalPolicy: fatalAlways}, title: "Не удалось получить объявление версии", allowedScopes: []model.DiagnosticScope{model.ScopeLine}, message: "Не удалось получить объявление версии"},
		{registryCase: registryCase{name: "P014", code: P014, scope: model.ScopeElement, fatalPolicy: fatalAlways}, title: "Версия не поддерживается", allowedScopes: []model.DiagnosticScope{model.ScopeElement}, details: Details{Received: "2.0", Supported: []string{"1.2", "1.3"}}, message: "Версия 2.0 не поддерживается; поддерживаемые версии: 1.2, 1.3"},
		{registryCase: registryCase{name: "IO001", code: IO001, scope: model.ScopeDocument, fatalPolicy: fatalAlways}, title: "Невозможно прочитать исходный файл", allowedScopes: []model.DiagnosticScope{model.ScopeDocument}, details: Details{Path: "/lesson.txt", Cause: readFailure}, message: "Невозможно прочитать исходный файл /lesson.txt: permission denied"},
	}

	definitions := definitions()
	if len(definitions) != len(tests) {
		t.Fatalf("definitions length = %d, want %d", len(definitions), len(tests))
	}
	wantCodes := make([]Code, len(tests))
	for position, test := range tests {
		wantCodes[position] = test.code
		t.Run(test.name, func(t *testing.T) {
			definition := definitions[position]
			if definition.code != test.code || definition.title != test.title || definition.fatal != test.fatalPolicy {
				t.Fatalf("definition = %#v, want code %q, title %q, fatal %v", definition, test.code, test.title, test.fatalPolicy)
			}
			if !reflect.DeepEqual(definition.scopes, test.allowedScopes) {
				t.Fatalf("scopes = %v, want %v", definition.scopes, test.allowedScopes)
			}
			if got := definition.message(test.details); got != test.message {
				t.Fatalf("message() = %q, want %q", got, test.message)
			}
			lookedUp, ok := lookup(test.code)
			if !ok || lookedUp.code != definition.code || lookedUp.title != definition.title {
				t.Fatalf("lookup(%q) = %#v, %v", test.code, lookedUp, ok)
			}
			for _, scope := range []model.DiagnosticScope{model.ScopeElement, model.ScopeLine, model.ScopeBlock, model.ScopeDocument, model.DiagnosticScope("unknown")} {
				want := containsScope(test.allowedScopes, scope)
				if got := acceptsScope(definition, scope); got != want {
					t.Errorf("acceptsScope(%q) = %v, want %v", scope, got, want)
				}
			}
			if err := validateFatal(test.fatalPolicy, true); err != nil {
				t.Errorf("validateFatal(policy, true) error: %v", err)
			}
			if err := validateFatal(test.fatalPolicy, false); (err == nil) != (test.fatalPolicy == fatalOnAmbiguity) {
				t.Errorf("validateFatal(policy, false) error = %v", err)
			}
		})
	}
	if got := Codes(); !reflect.DeepEqual(got, wantCodes) {
		t.Fatalf("Codes() = %v, want %v", got, wantCodes)
	}
	first := Codes()
	first[0] = IO001
	if reflect.DeepEqual(first, Codes()) {
		t.Fatal("Codes() exposes registry storage")
	}
	if _, ok := lookup(Code("P999")); ok {
		t.Fatal("lookup(unknown) returned ok=true")
	}
	if err := validateFatal(fatalPolicy(99), true); err == nil {
		t.Fatal("validateFatal(unknown policy) error = nil")
	}
}

// containsScope сообщает, содержится ли scope в разрешённом тестовом наборе.
func containsScope(scopes []model.DiagnosticScope, scope model.DiagnosticScope) bool {
	for _, candidate := range scopes {
		if candidate == scope {
			return true
		}
	}
	return false
}
