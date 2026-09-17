package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
)

// tagName задаёт закрытый словарь канонических имён тегов DSL v1.2.
type tagName string

// supportedTagCount содержит число канонических тегов в закрытом словаре DSL v1.2.
const supportedTagCount = 34

const (
	// tagDSLVersion обозначает тег dsl-version.
	tagDSLVersion tagName = "dsl-version"
	// tagDocumentID обозначает тег document-id.
	tagDocumentID tagName = "document-id"
	// tagSection обозначает тег section.
	tagSection tagName = "section"
	// tagOrder обозначает тег order.
	tagOrder tagName = "order"
	// tagResourceDir обозначает тег resource-dir.
	tagResourceDir tagName = "resource-dir"
	// tagHeader обозначает тег header.
	tagHeader tagName = "header"
	// tagTask обозначает тег task.
	tagTask tagName = "task"
	// tagEndTask обозначает тег endtask.
	tagEndTask tagName = "endtask"
	// tagStep обозначает тег step.
	tagStep tagName = "step"
	// tagSpeaking обозначает тег speaking.
	tagSpeaking tagName = "speaking"
	// tagNewPage обозначает тег newpage.
	tagNewPage tagName = "newpage"
	// tagEditor обозначает тег editor.
	tagEditor tagName = "editor"
	// tagMedia обозначает тег media.
	tagMedia tagName = "media"
	// tagExample обозначает тег example.
	tagExample tagName = "example"
	// tagWordlist обозначает тег wordlist.
	tagWordlist tagName = "wordlist"
	// tagTable обозначает тег table.
	tagTable tagName = "table"
	// tagScript обозначает тег script.
	tagScript tagName = "script"
	// tagText обозначает тег text.
	tagText tagName = "text"
	// tagKey обозначает тег key.
	tagKey tagName = "key"
	// tagInstruction обозначает тег instr.
	tagInstruction tagName = "instr"
	// tagNote обозначает тег note.
	tagNote tagName = "note"
	// tagAlt обозначает тег alt.
	tagAlt tagName = "alt"
	// tagHint обозначает тег hint.
	tagHint tagName = "hint"
	// tagFragment обозначает тег fragment.
	tagFragment tagName = "fragment"
	// tagInclude обозначает тег include.
	tagInclude tagName = "include"
	// tagAnswer обозначает тег answer.
	tagAnswer tagName = "answer"
	// tagQuestion обозначает тег question.
	tagQuestion tagName = "question"
	// tagMultifill обозначает тег multifill.
	tagMultifill tagName = "multifill"
	// tagChoice обозначает тег choice.
	tagChoice tagName = "choice"
	// tagMultichoice обозначает тег multichoice.
	tagMultichoice tagName = "multichoice"
	// tagMatching обозначает тег matching.
	tagMatching tagName = "matching"
	// tagOrdering обозначает тег ordering.
	tagOrdering tagName = "ordering"
	// tagVariants обозначает тег variants.
	tagVariants tagName = "variants"
	// tagVariant обозначает тег variant.
	tagVariant tagName = "variant"
)

// syntaxKind задаёт структуру операндов одной формы объявления.
type syntaxKind int

const (
	// syntaxNone обозначает форму без операндов.
	syntaxNone syntaxKind = iota
	// syntaxRequiredToken обозначает один обязательный отделённый операнд.
	syntaxRequiredToken
	// syntaxRequiredFree обозначает одно обязательное свободное значение до конца формы.
	syntaxRequiredFree
	// syntaxOptionalFree обозначает необязательное свободное значение до конца формы.
	syntaxOptionalFree
	// syntaxMedia обозначает обязательные TYPE и SOURCE объявления media.
	syntaxMedia
	// syntaxResourcePaths обозначает список путей resource-dir.
	syntaxResourcePaths
)

// contentPolicy задаёт контекстные правила экранирования и специальных последовательностей.
type contentPolicy int

const (
	// contentPlain обрабатывает обычное текстовое значение DSL.
	contentPlain contentPolicy = iota
	// contentWordlist дополнительно обрабатывает структурную точку с запятой.
	contentWordlist
	// contentExample дополнительно защищает строку-разделитель example.
	contentExample
	// contentTable обрабатывает специальные символы строки table без выделения ячеек.
	contentTable
	// contentHTMLText обрабатывает литеральные угловые скобки содержимого text.
	contentHTMLText
	// contentMultifill защищает плейсхолдеры без создания placeholder-элементов.
	contentMultifill
	// contentEditor сохраняет похожие на DSL конструкции как редакторский текст.
	contentEditor
	// contentMediaSource обрабатывает кавычки и обратную косую черту SOURCE.
	contentMediaSource
	// contentResourcePath сохраняет обратную косую черту пути буквально.
	contentResourcePath
)

// formSpec задаёт операнды и тип значения одной структурной формы тега.
type formSpec struct {
	// syntax содержит структуру операндов формы.
	syntax syntaxKind
	// element содержит тип единственного обычного значения либо пустую строку для специальной формы.
	element model.ElementType
	// policy содержит правила снятия экранирования значения.
	policy contentPolicy
}

// parentRule задаёт версионное правило построения упорядоченных кандидатов родителя.
type parentRule int

const (
	// parentDocument задаёт корневую строку документа.
	parentDocument parentRule = iota
	// parentContent задаёт текущий блок, внутренний step, task, variant или шаг документа.
	parentContent
	// parentTask задаёт активное задание.
	parentTask
	// parentStep задаёт текущий шаг документа.
	parentStep
	// parentVariants задаёт внешний блок variants.
	parentVariants
	// parentVariantOrStep задаёт текущую ветвь либо шаг документа.
	parentVariantOrStep
)

// blockSpec задаёт контекст и лексический режим открываемого блока.
type blockSpec struct {
	// context содержит контекст содержимого блока.
	context grammar.Context
	// mode содержит способ лексического восприятия содержимого блока.
	mode grammar.ContentMode
}

// tagSpec содержит неизменяемое декларативное описание одного тега DSL v1.2.
type tagSpec struct {
	// name содержит каноническое имя тега без символа @.
	name tagName
	// forms содержит разрешённые структурные формы.
	forms grammar.TagForm
	// line содержит синтаксис однострочной формы.
	line formSpec
	// blockForm содержит синтаксис блочной формы перед открывающей скобкой.
	blockForm formSpec
	// lineParent содержит правило родителя однострочной формы.
	lineParent parentRule
	// blockParent содержит правило родителя блочной формы.
	blockParent parentRule
	// block содержит режим открываемого блока.
	block blockSpec
	// transition содержит изменение неблочного состояния.
	transition grammar.Transition
}

// tagSpecs возвращает отдельный упорядоченный набор описаний всех тегов DSL v1.2.
func tagSpecs() []tagSpec { panic("TODO") }

// findTagSpec возвращает описание тега и true либо нулевое значение и false для неизвестного имени.
func findTagSpec(name tagName) (tagSpec, bool) { panic("TODO") }
