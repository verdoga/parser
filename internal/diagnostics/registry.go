package diagnostics

import "dslparser/internal/model"

// fatalPolicy задаёт допустимость флага fatal для кода.
type fatalPolicy int

const (
	// fatalAlways требует fatal=true.
	fatalAlways fatalPolicy = iota
	// fatalOnAmbiguity допускает false при восстановлении и true при неоднозначном состоянии.
	fatalOnAmbiguity
)

// messageFunc формирует нормативное сообщение из типизированных деталей.
type messageFunc func(details Details) string

// definition содержит неизменяемый контракт одного кода реестра.
type definition struct {
	// code содержит стабильный машинный код.
	code Code
	// title содержит нормативное русское название ошибки из реестра.
	title string
	// scopes содержит допустимые области диагностики.
	scopes []model.DiagnosticScope
	// fatal содержит политику фатальности.
	fatal fatalPolicy
	// message формирует человекочитаемое сообщение.
	message messageFunc
}

// definitions возвращает новый регистрационный список всех кодов в нормативном порядке.
func definitions() []definition { panic("TODO") }

// lookup возвращает запись реестра и true либо нулевую запись и false для неизвестного кода.
func lookup(code Code) (definition, bool) { panic("TODO") }

// Codes возвращает отдельную копию кодов в нормативном порядке P001–P014, IO001.
func Codes() []Code { panic("TODO") }

// acceptsScope сообщает true для разрешённой области и false для запрещённой.
func acceptsScope(definition definition, scope model.DiagnosticScope) bool { panic("TODO") }

// validateFatal проверяет флаг fatal по политике кода.
func validateFatal(policy fatalPolicy, fatal bool) error { panic("TODO") }

// messageP001 формирует сообщение о недопустимой последовательности UTF-8.
func messageP001(details Details) string { panic("TODO") }

// messageP002 формирует сообщение об одиночном CR.
func messageP002(details Details) string { panic("TODO") }

// messageP003 формирует сообщение о неизвестной DSL-конструкции.
func messageP003(details Details) string { panic("TODO") }

// messageP004 формирует сообщение о нарушенном обязательном разделителе.
func messageP004(details Details) string { panic("TODO") }

// messageP005 формирует сообщение о неподдерживаемой форме известного тега.
func messageP005(details Details) string { panic("TODO") }

// messageP006 формирует сообщение об отсутствующем структурном параметре.
func messageP006(details Details) string { panic("TODO") }

// messageP007 формирует сообщение о лишнем фрагменте объявления.
func messageP007(details Details) string { panic("TODO") }

// messageP008 формирует сообщение о неправильной строке открытия блока.
func messageP008(details Details) string { panic("TODO") }

// messageP009 формирует сообщение о закрытии при пустом стеке блоков.
func messageP009(details Details) string { panic("TODO") }

// messageP010 формирует сообщение о совмещённой с содержимым закрывающей скобке.
func messageP010(details Details) string { panic("TODO") }

// messageP011 формирует сообщение о незакрытом блоке с именем открывающего тега.
func messageP011(details Details) string { panic("TODO") }

// messageP012 формирует сообщение о неэкранированной фигурной скобке в тексте.
func messageP012(details Details) string { panic("TODO") }

// messageP013 формирует сообщение о невозможности выбрать версию по первой строке.
func messageP013(details Details) string { panic("TODO") }

// messageP014 формирует сообщение с полученной и поддерживаемыми версиями.
func messageP014(details Details) string { panic("TODO") }

// messageIO001 формирует сообщение о невозможности прочитать конкретный источник.
func messageIO001(details Details) string { panic("TODO") }
