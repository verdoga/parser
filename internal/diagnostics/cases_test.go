package diagnostics

import "dslparser/internal/model"

// registryCase описывает табличную проверку записи единого реестра.
type registryCase struct {
	// name содержит имя подтеста.
	name string
	// code содержит проверяемый код.
	code Code
	// scope содержит ожидаемую допустимую область.
	scope model.DiagnosticScope
	// fatalPolicy содержит ожидаемую политику фатальности.
	fatalPolicy fatalPolicy
}

// builderCase описывает проверку создания, порядка или дедупликации.
type builderCase struct {
	// name содержит имя подтеста.
	name string
	// request содержит входной запрос накопителя.
	request Request
	// created содержит ожидаемый признак создания новой диагностики.
	created bool
}

// messageCase описывает проверку нормативного сообщения с деталями.
type messageCase struct {
	// name содержит имя подтеста.
	name string
	// code содержит код форматируемого сообщения.
	code Code
	// details содержит переменные части сообщения.
	details Details
	// want содержит ожидаемое нормативное сообщение.
	want string
}
