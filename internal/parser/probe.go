package parser

import "dslparser/internal/model"

// versionProbe содержит результат синтаксического чтения первой строки.
type versionProbe struct {
	// value содержит выделенную версию либо nil.
	value *string
	// element содержит выделенный элемент версии либо nil.
	element *model.Element
	// problem содержит P013 либо nil для корректной формы.
	problem *diagnosticDraft
}

// probeVersion разбирает только нормативную форму версии в первой физической строке.
func probeVersion(line physicalLine) versionProbe { panic("TODO") }

// probeDocumentID возвращает единственный однозначно выделенный ID либо nil.
func probeDocumentID(lines []physicalLine, selected Grammar) *string { panic("TODO") }
