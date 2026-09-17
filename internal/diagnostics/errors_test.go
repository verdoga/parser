package diagnostics

import "testing"

// TestRequestError проверяет человекочитаемое представление устойчивых полей ошибки.
func TestRequestError(t *testing.T) {
	tests := []struct {
		name  string
		error RequestError
		want  string
	}{
		{name: "поле и правило", error: RequestError{Field: "location", Rule: "диапазон обязателен"}, want: "location: диапазон обязателен"},
		{name: "пустое поле", error: RequestError{Rule: "нарушен контракт"}, want: "нарушен контракт"},
		{name: "пустое правило", error: RequestError{Field: "code"}, want: "code"},
		{name: "нулевое значение", error: RequestError{}, want: "недопустимый запрос диагностики"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.error.Error(); got != test.want {
				t.Fatalf("Error() = %q, want %q", got, test.want)
			}
		})
	}
}
