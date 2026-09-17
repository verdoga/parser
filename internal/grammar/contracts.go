package grammar

import "dslparser/internal/model"

// Grammar задаёт реальную границу подмены правил одной версии DSL.
type Grammar interface {
	// Version возвращает точное поддерживаемое значение версии.
	Version() string
	// Classify классифицирует строку в заданном контексте без изменения состояния parser.
	Classify(request GrammarRequest) GrammarDecision
}

// TagForm задаёт структурную форму объявления тега.
type TagForm int

const (
	// TagFormLine обозначает однострочное объявление без открывающей границы блока.
	TagFormLine TagForm = 1 << iota
	// TagFormBlock обозначает объявление, открывающее фигурный блок.
	TagFormBlock
)

// Allows сообщает, разрешена ли одиночная форма candidate набором форм f.
// Результат true означает, что форма разрешена; false означает, что форма не входит в набор.
func (f TagForm) Allows(candidate TagForm) bool { panic("TODO") }

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
// Parser выбирает Parent по состоянию перед строкой, затем применяет Transition и после него Block.
// Для BlockClose родителем остаётся закрываемая вершина стека, а само удаление выполняется последним.
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

// Transition задаёт завершение и открытие неблочных областей после выбора родителя строки.
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

// ParentDecision задаёт правило определения логического родителя до перехода состояния строки.
type ParentDecision struct {
	// Kinds содержит источники родителя в порядке приоритета.
	// Parser выбирает первую существующую область; пустой срез задаёт корневую строку.
	Kinds []ParentKind
}

// Parents создаёт правило выбора первого существующего логического родителя.
// Возвращаемое значение владеет отдельной копией переданного порядка.
func Parents(primary ParentKind, fallback ...ParentKind) ParentDecision { panic("TODO") }

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
