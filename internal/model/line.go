package model

// LineType задаёт закрытое перечисление типов физических строк.
type LineType string

const (
	// LineTag обозначает однострочное объявление тега.
	LineTag LineType = "tag"
	// LineHeading обозначает заголовок DSL.
	LineHeading LineType = "heading"
	// LineContent обозначает строку содержимого.
	LineContent LineType = "content"
	// LineBlockStart обозначает объявление с открывающей фигурной скобкой.
	LineBlockStart LineType = "block-start"
	// LineBlockEnd обозначает отдельную закрывающую фигурную скобку.
	LineBlockEnd LineType = "block-end"
	// LineSeparator обозначает разрешённую отдельную строку разделителя.
	LineSeparator LineType = "separator"
	// LineBlank обозначает игнорируемую пустую строку.
	LineBlank LineType = "blank"
	// LineInvalid обозначает строку, которую нельзя классифицировать.
	LineInvalid LineType = "invalid"
)

// LineEnding задаёт закрытое перечисление сохранённых окончаний строки.
type LineEnding string

const (
	// EOLNone обозначает EOF без перевода строки.
	EOLNone LineEnding = ""
	// EOLLF обозначает перевод строки LF.
	EOLLF LineEnding = "\n"
	// EOLCRLF обозначает перевод строки CRLF.
	EOLCRLF LineEnding = "\r\n"
)

// ElementType задаёт закрытое перечисление типов элементов строки.
type ElementType string

const (
	// ElementTag обозначает каноническое имя тега.
	ElementTag ElementType = "tag"
	// ElementHeadingLevel обозначает синтаксический уровень заголовка.
	ElementHeadingLevel ElementType = "heading-level"
	// ElementTitle обозначает свободный текст названия.
	ElementTitle ElementType = "title"
	// ElementContent обозначает текстовое содержимое.
	ElementContent ElementType = "content"
	// ElementIdentifier обозначает идентификатор DSL.
	ElementIdentifier ElementType = "identifier"
	// ElementName обозначает свободное имя.
	ElementName ElementType = "name"
	// ElementVersion обозначает прочитанную версию DSL.
	ElementVersion ElementType = "version"
	// ElementNumber обозначает исходную числовую запись.
	ElementNumber ElementType = "number"
	// ElementMediaType обозначает тип media.
	ElementMediaType ElementType = "media-type"
	// ElementSource обозначает источник media.
	ElementSource ElementType = "source"
	// ElementResourcePath обозначает один путь resource-dir.
	ElementResourcePath ElementType = "resource-path"
	// ElementPlaceholder обозначает общий контракт плейсхолдера, который parser не создаёт.
	ElementPlaceholder ElementType = "placeholder"
	// ElementUnparsed обозначает ошибочный фрагмент с element-диагностикой.
	ElementUnparsed ElementType = "unparsed"
	// ElementBlockOpen обозначает структурную открывающую скобку.
	ElementBlockOpen ElementType = "block-open"
	// ElementBlockClose обозначает структурную закрывающую скобку.
	ElementBlockClose ElementType = "block-close"
)

// Line содержит полное построчное представление одной физической строки.
type Line struct {
	// Number содержит уникальный номер физической строки от единицы.
	Number int `json:"line"`
	// Type содержит классификацию строки.
	Type LineType `json:"type"`
	// NestingLevel содержит логическую глубину, согласованную с родителем.
	NestingLevel int `json:"nestingLevel"`
	// ParentLine содержит номер предшествующей родительской строки либо nil.
	ParentLine *int `json:"parentLine"`
	// Raw содержит точную строку без BOM и окончания строки.
	Raw string `json:"raw"`
	// EOL содержит точное поддерживаемое окончание строки.
	EOL LineEnding `json:"eol"`
	// HasErrors равен true при собственной error-диагностике строки и false иначе.
	HasErrors bool `json:"hasErrors"`
	// Elements содержит непересекающиеся элементы в порядке диапазонов.
	Elements []Element `json:"elements"`
}

// Element содержит типизированный непересекающийся фрагмент строки.
type Element struct {
	// Type содержит семантический тип фрагмента.
	Type ElementType `json:"type"`
	// Raw содержит точный фрагмент строки по диапазону Unicode-колонок.
	Raw string `json:"raw"`
	// Value содержит семантическое значение либо допустимый для типа nil.
	Value *string `json:"value"`
	// Start содержит включительную Unicode-колонку от единицы.
	Start int `json:"start"`
	// End содержит исключающую Unicode-колонку, не меньшую Start.
	End int `json:"end"`
	// ErrorIDs содержит уникальные ссылки на element-диагностики в нормативном порядке.
	ErrorIDs []string `json:"errorIds"`
}
