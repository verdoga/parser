package parser

import "dslparser/internal/model"

// Registry предоставляет неизменяемый набор реализаций поддерживаемых версий DSL.
type Registry interface {
	// Lookup возвращает грамматику версии и true либо nil и false для неподдерживаемой версии.
	Lookup(version string) (Grammar, bool)
	// Versions возвращает отдельную копию поддерживаемых версий в детерминированном порядке.
	Versions() []string
}

// Grammar задаёт реальную границу подмены правил одной версии DSL.
type Grammar interface {
	// Version возвращает точное поддерживаемое значение версии.
	Version() string
	// Classify классифицирует строку в заданном контексте без изменения состояния parser.
	Classify(request GrammarRequest) GrammarDecision
}

// ContentMode задаёт способ лексического восприятия содержимого текущего блока.
type ContentMode int

const (
	// ContentDSL распознаёт структурные конструкции DSL.
	ContentDSL ContentMode = iota
	// ContentOpaque сохраняет строки как непрозрачное содержимое.
	ContentOpaque
	// ContentMultifill сохраняет содержимое и защищает допустимые скобки ответа плейсхолдера.
	ContentMultifill
	// ContentEditor сохраняет похожие на DSL конструкции как редакторский текст.
	ContentEditor
)

// EscapeMode задаёт контекстное снятие экранирования.
type EscapeMode int

const (
	// EscapeText обрабатывает синтаксически значимое экранирование обычного текста.
	EscapeText EscapeMode = iota
	// EscapeMedia обрабатывает кавычки и обратную косую черту источника media.
	EscapeMedia
	// EscapeResourcePath сохраняет буквальную обратную косую черту пути ресурса.
	EscapeResourcePath
	// EscapeNone не снимает экранирование.
	EscapeNone
)

// Context содержит устойчивое имя грамматического контекста версии.
type Context string

// GrammarRequest содержит наблюдаемую строку и состояние, необходимые грамматике версии.
type GrammarRequest struct {
	// Raw содержит точную физическую строку без EOL.
	Raw string
	// Trimmed содержит строку после удаления крайних U+0020 и TAB.
	Trimmed string
	// Line содержит номер физической строки.
	Line int
	// Context содержит текущий контекст grammar.
	Context Context
	// ContentMode содержит текущий режим содержимого.
	ContentMode ContentMode
}

// GrammarDecision содержит полное декларативное решение grammar для одной строки.
type GrammarDecision struct {
	// LineType содержит классификацию строки.
	LineType model.LineType
	// Elements содержит описания выделенных элементов в порядке диапазонов.
	Elements []ElementDecision
	// Transition содержит изменение неблочного состояния.
	Transition Transition
	// Parent содержит правило выбора логического родителя.
	Parent ParentDecision
	// Block содержит необязательное действие со стеком блоков.
	Block BlockDecision
	// Recovery содержит возможность однозначного продолжения после ошибки.
	Recovery RecoveryDecision
	// Problems содержит обнаруженные грамматикой проблемы в порядке обнаружения.
	Problems []Problem
}

// ElementDecision описывает один выделенный grammar элемент через байтовые границы исходной строки.
type ElementDecision struct {
	// Type содержит тип будущего model.Element.
	Type model.ElementType
	// Value содержит семантическое значение либо nil.
	Value *string
	// ByteStart содержит включительное смещение начала в UTF-8 строке.
	ByteStart int
	// ByteEnd содержит исключающее смещение конца в UTF-8 строке.
	ByteEnd int
}

// Transition задаёт завершение и открытие неблочных областей.
type Transition struct {
	// CloseTask завершает активное задание, если оно существует.
	CloseTask bool
	// OpenTask открывает задание на текущей строке.
	OpenTask bool
	// SetStep устанавливает текущий шаг документа на текущую строку.
	SetStep bool
	// SetInnerStep устанавливает внутренний шаг задания на текущую строку.
	SetInnerStep bool
	// OpenVariants устанавливает блок variants на текущую строку.
	OpenVariants bool
	// SetVariant устанавливает текущую ветвь на текущую строку.
	SetVariant bool
	// ClearVariants завершает область variants и текущую ветвь.
	ClearVariants bool
}

// ParentKind задаёт источник логического родителя строки.
type ParentKind int

const (
	// ParentRoot задаёт отсутствие родителя.
	ParentRoot ParentKind = iota
	// ParentBlock задаёт верхний фигурный блок.
	ParentBlock
	// ParentTask задаёт активное задание.
	ParentTask
	// ParentStep задаёт текущий шаг документа.
	ParentStep
	// ParentInnerStep задаёт внутренний шаг задания.
	ParentInnerStep
	// ParentVariants задаёт блок variants.
	ParentVariants
	// ParentVariant задаёт текущую ветвь variant.
	ParentVariant
)

// ParentDecision задаёт правило определения логического родителя.
type ParentDecision struct {
	// Kind содержит источник родителя.
	Kind ParentKind
}

// BlockAction задаёт изменение стека фигурных блоков.
type BlockAction int

const (
	// BlockUnchanged не изменяет стек.
	BlockUnchanged BlockAction = iota
	// BlockOpen добавляет текущую строку в стек.
	BlockOpen
	// BlockClose удаляет верхний блок после сопоставления.
	BlockClose
)

// BlockDecision описывает действие со стеком и контекст открываемого блока.
type BlockDecision struct {
	// Action содержит действие со стеком.
	Action BlockAction
	// Tag содержит каноническое имя открываемого тега.
	Tag string
	// Context содержит контекст содержимого открываемого блока.
	Context Context
	// ContentMode содержит режим содержимого открываемого блока.
	ContentMode ContentMode
}

// RecoveryDecision задаёт однозначность состояния после ошибки строки.
type RecoveryDecision struct {
	// Continue равен true при однозначном продолжении и false при необходимости остановки.
	Continue bool
}

// Problem описывает обнаруженную grammar проблему до назначения диагностического ID.
type Problem struct {
	// Code содержит один из кодов P003–P012.
	Code string
	// Scope содержит область будущей диагностики.
	Scope model.DiagnosticScope
	// Fatal равен true при неоднозначном продолжении и false при восстановлении.
	Fatal bool
	// ByteStart содержит включительное байтовое смещение основного диапазона.
	ByteStart int
	// ByteEnd содержит исключающее байтовое смещение основного диапазона.
	ByteEnd int
	// Element содержит индекс связанного ElementDecision либо -1 без связи.
	Element int
}
