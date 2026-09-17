package parser

import "dslparser/internal/model"

// metadataCollector накапливает все синтаксически выделенные кандидаты метаданных.
type metadataCollector struct {
	// documentIDs содержит объявления document-id в исходном порядке.
	documentIDs []locatedValue
	// titles содержит заголовки первого уровня в исходном порядке.
	titles []locatedValue
	// subtitles содержит заголовки второго уровня в исходном порядке.
	subtitles []locatedValue
	// sections содержит объявления section в исходном порядке.
	sections []locatedValue
	// orders содержит объявления order в исходном порядке.
	orders []locatedValue
	// resourceDirs содержит пути ресурсов в порядке объявлений.
	resourceDirs []locatedValue
}

// locatedValue связывает сохранённое значение с физической строкой объявления.
type locatedValue struct {
	// value содержит значение с сохранённым исходным регистром.
	value string
	// line содержит номер физической строки объявления.
	line int
}

// observe учитывает метаданные уже построенной строки, не изменяя её элементы.
func (c *metadataCollector) observe(line model.Line) { panic("TODO") }

// metadata возвращает однозначные поля и различает недоступный и пустой список ресурсов.
func (c metadataCollector) metadata(available bool) model.Metadata { panic("TODO") }

// uniqueValue возвращает отдельную копию единственного значения либо nil при ином количестве.
func uniqueValue(values []locatedValue) *string { panic("TODO") }
