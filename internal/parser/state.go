package parser

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// parseState содержит независимые состояния фигурных и неблочных областей.
type parseState struct {
	// blocks содержит стек фигурных блоков от внешнего к внутреннему.
	blocks []blockFrame
	// currentStep содержит строку текущего шага документа либо ноль.
	currentStep int
	// activeTask содержит строку активного задания либо ноль.
	activeTask int
	// innerStep содержит строку внутреннего шага задания либо ноль.
	innerStep int
	// variants содержит строку активного блока variants либо ноль.
	variants int
	// variant содержит строку текущей ветви variant либо ноль.
	variant int
	// stopped сообщает, был ли структурный разбор фатально остановлен.
	stopped bool
}

// blockFrame содержит открывающую строку и режим одного фигурного блока.
type blockFrame struct {
	// tag содержит каноническое имя открывающего тега.
	tag string
	// openingLine содержит номер строки block-start.
	openingLine int
	// opening содержит диапазон структурной открывающей скобки.
	opening model.Range
	// mode содержит режим распознавания содержимого.
	mode grammar.ContentMode
	// context содержит грамматический контекст содержимого.
	context grammar.Context
}

// newParseState создаёт независимое пустое состояние одного разбора.
func newParseState() parseState { panic("TODO") }

// openBlock добавляет однозначно открытый блок на вершину стека.
func (s *parseState) openBlock(frame blockFrame) { panic("TODO") }

// closeBlock удаляет и возвращает верхний блок.
func (s *parseState) closeBlock() blockFrame { panic("TODO") }

// topBlock возвращает верхний блок и true либо нулевое значение и false для пустого стека.
func (s parseState) topBlock() (blockFrame, bool) { panic("TODO") }

// applyTransition применяет однозначный переход неблочных областей к текущей строке.
func (s *parseState) applyTransition(transition grammar.Transition, line int) { panic("TODO") }

// parentFor возвращает номер первого существующего логического родителя либо ноль для корневой строки.
func (s parseState) parentFor(decision grammar.ParentDecision) int { panic("TODO") }
