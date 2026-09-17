package report

import "dslparser/internal/model"

// Компиляционные проверки фиксируют минимальный контракт сборки и сериализации отчёта.
var (
	_ = Input{
		Document:    model.Document{},
		Processing:  model.Processing{},
		Lines:       []model.Line{},
		Diagnostics: []model.Diagnostic{},
	}
	_ func(Input) model.Result           = Build
	_ func(model.Result) ([]byte, error) = Marshal
)
