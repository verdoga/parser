package v1_2

import "dslparser/internal/grammar"

// Компиляционные проверки фиксируют внутренние обработчики форм объявлений.
var (
	_ declarationParser                                    = parseNoArguments
	_ declarationParser                                    = parseRequiredArgument
	_ declarationParser                                    = parseFreeValue
	_ declarationParser                                    = parseMedia
	_ declarationParser                                    = parseResourceDirs
	_ declarationParser                                    = parseMultifill
	_ func(grammar.GrammarRequest) grammar.GrammarDecision = classifyHeading
	_ func(grammar.GrammarRequest) grammar.GrammarDecision = classifyBlockBoundary
	_ func(grammar.GrammarRequest) grammar.GrammarDecision = classifyDeclaration
	_ func(grammar.GrammarRequest) grammar.GrammarDecision = classifyContent
	_ func(grammar.GrammarRequest) grammar.GrammarDecision = classifySeparator
)

// Компиляционные проверки фиксируют полный закрытый словарь из 34 тегов DSL v1.2.
var (
	_ int     = supportedTagCount
	_ tagName = tagDSLVersion
	_ tagName = tagDocumentID
	_ tagName = tagSection
	_ tagName = tagOrder
	_ tagName = tagResourceDir
	_ tagName = tagHeader
	_ tagName = tagTask
	_ tagName = tagEndTask
	_ tagName = tagStep
	_ tagName = tagSpeaking
	_ tagName = tagNewPage
	_ tagName = tagEditor
	_ tagName = tagMedia
	_ tagName = tagExample
	_ tagName = tagWordlist
	_ tagName = tagTable
	_ tagName = tagScript
	_ tagName = tagText
	_ tagName = tagKey
	_ tagName = tagInstruction
	_ tagName = tagNote
	_ tagName = tagAlt
	_ tagName = tagHint
	_ tagName = tagFragment
	_ tagName = tagInclude
	_ tagName = tagAnswer
	_ tagName = tagQuestion
	_ tagName = tagMultifill
	_ tagName = tagChoice
	_ tagName = tagMultichoice
	_ tagName = tagMatching
	_ tagName = tagOrdering
	_ tagName = tagVariants
	_ tagName = tagVariant
)
