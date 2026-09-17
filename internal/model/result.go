package model

// FormatVersion содержит поддерживаемую версию JSON-контракта, а не версию DSL.
const FormatVersion = "1.0"

// Result содержит полный результат обработки, готовый к проверке и сериализации.
// Экспортированные поля являются намеренным транспортным контрактом для encoding/json.
type Result struct {
	// FormatVersion содержит версию JSON-контракта.
	FormatVersion string `json:"formatVersion"`
	// Document содержит сведения об исходном документе.
	Document Document `json:"document"`
	// Processing содержит запуски инструментов в порядке их начала.
	Processing []Processing `json:"processing"`
	// Lines содержит физические строки в исходном порядке.
	Lines []Line `json:"lines"`
	// Diagnostics содержит диагностики в нормативном порядке.
	Diagnostics []Diagnostic `json:"diagnostics"`
}
