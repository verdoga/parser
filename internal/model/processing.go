package model

// Processing описывает один фактически начатый запуск инструмента.
type Processing struct {
	// ID содержит непустой уникальный в пределах результата идентификатор.
	ID string `json:"id"`
	// Tool содержит стабильное машинное имя инструмента.
	Tool string `json:"tool"`
	// Version содержит версию инструмента либо nil, если она не сообщена.
	Version *string `json:"version"`
	// StartedAt содержит время RFC 3339 с часовым поясом либо nil.
	StartedAt *string `json:"startedAt"`
	// DurationMS содержит неотрицательную длительность в миллисекундах либо nil.
	DurationMS *int `json:"durationMs"`
}
