package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"testing"
)

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

// TestTagForms проверяет закрытый словарь и все поддерживаемые формы 34 тегов.
func TestTagForms(t *testing.T) {
	tests := []struct {
		name        string
		forms       grammar.TagForm
		lineSyntax  syntaxKind
		blockSyntax syntaxKind
	}{
		{name: "dsl-version", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredToken},
		{name: "document-id", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredToken},
		{name: "section", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredFree},
		{name: "order", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredToken},
		{name: "resource-dir", forms: grammar.TagFormLine, lineSyntax: syntaxResourcePaths},
		{name: "header", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredFree},
		{name: "task", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredToken},
		{name: "endtask", forms: grammar.TagFormLine, lineSyntax: syntaxNone},
		{name: "step", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredFree},
		{name: "speaking", forms: grammar.TagFormLine, lineSyntax: syntaxNone},
		{name: "newpage", forms: grammar.TagFormLine, lineSyntax: syntaxNone},
		{name: "editor", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxNone},
		{name: "media", forms: grammar.TagFormLine, lineSyntax: syntaxMedia},
		{name: "example", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxOptionalFree},
		{name: "wordlist", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxOptionalFree},
		{name: "table", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "script", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "text", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "key", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "instr", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxNone},
		{name: "note", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxOptionalFree},
		{name: "alt", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxOptionalFree},
		{name: "hint", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxNone},
		{name: "fragment", forms: grammar.TagFormBlock, blockSyntax: syntaxRequiredToken},
		{name: "include", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredToken},
		{name: "answer", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxRequiredFree, blockSyntax: syntaxNone},
		{name: "question", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredFree},
		{name: "multifill", forms: grammar.TagFormLine | grammar.TagFormBlock, lineSyntax: syntaxOptionalFree, blockSyntax: syntaxOptionalFree},
		{name: "choice", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "multichoice", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "matching", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "ordering", forms: grammar.TagFormBlock, blockSyntax: syntaxOptionalFree},
		{name: "variants", forms: grammar.TagFormBlock, blockSyntax: syntaxNone},
		{name: "variant", forms: grammar.TagFormLine, lineSyntax: syntaxRequiredFree},
	}

	specs := tagSpecs()
	if len(specs) != supportedTagCount || len(tests) != supportedTagCount {
		t.Fatalf("tagSpecs length = %d, test length = %d, want %d", len(specs), len(tests), supportedTagCount)
	}
	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := specs[index]
			if string(spec.name) != test.name || spec.forms != test.forms || spec.line.syntax != test.lineSyntax || spec.blockForm.syntax != test.blockSyntax {
				t.Fatalf("tagSpecs[%d] = %#v", index, spec)
			}
			found, ok := findTagSpec(tagName(test.name))
			if !ok || found.name != spec.name || found.forms != spec.forms {
				t.Fatalf("findTagSpec(%q) = (%#v, %t)", test.name, found, ok)
			}
			if test.forms.Allows(grammar.TagFormLine) {
				line := validDeclaration(spec, grammar.TagFormLine)
				decision := New().Classify(request(line, contextRoot, grammar.ContentDSL))
				if decision.LineType != model.LineTag || len(decision.Problems) != 0 {
					t.Fatalf("valid line %q decision = %#v", line, decision)
				}
			}
			if test.forms.Allows(grammar.TagFormBlock) {
				line := validDeclaration(spec, grammar.TagFormBlock)
				decision := New().Classify(request(line, contextRoot, grammar.ContentDSL))
				if decision.LineType != model.LineBlockStart || len(decision.Problems) != 0 {
					t.Fatalf("valid block %q decision = %#v", line, decision)
				}
			}
		})
	}
	if _, ok := findTagSpec("Task"); ok {
		t.Fatal("findTagSpec accepted non-canonical case")
	}
	if _, ok := findTagSpec("unknown"); ok {
		t.Fatal("findTagSpec accepted unknown tag")
	}
	specs[0].name = "changed"
	if fresh := tagSpecs(); fresh[0].name != tagDSLVersion {
		t.Fatal("tagSpecs exposed mutable package state")
	}
}

// validDeclaration создаёт минимальное корректное объявление заданной формы.
func validDeclaration(spec tagSpec, form grammar.TagForm) string {
	result := "@" + string(spec.name)
	formSpec := spec.line
	if form == grammar.TagFormBlock {
		formSpec = spec.blockForm
	}
	switch formSpec.syntax {
	case syntaxRequiredToken, syntaxRequiredFree:
		result += " value"
	case syntaxMedia:
		result += " audio source.mp3"
	case syntaxResourcePaths:
		result += " assets"
	}
	if form == grammar.TagFormBlock {
		result += " {"
	}
	return result
}

// TestTagClassificationErrors проверяет регистр, U+0020, TAB и неверные части объявления.
func TestTagClassificationErrors(t *testing.T) {
	tests := []struct {
		name string
		line string
		kind grammar.ProblemKind
	}{
		{name: "unknown", line: "@unknown", kind: grammar.ProblemUnknownTag},
		{name: "case", line: "@Task id", kind: grammar.ProblemUnknownTag},
		{name: "tab separator", line: "@task\tid", kind: grammar.ProblemMissingSeparator},
		{name: "missing argument", line: "@task", kind: grammar.ProblemMissingArgument},
		{name: "extra no-argument content", line: "@newpage now", kind: grammar.ProblemExtraContent},
		{name: "unsupported line", line: "@table title", kind: grammar.ProblemUnsupportedForm},
		{name: "unsupported block", line: "@question text {", kind: grammar.ProblemUnsupportedForm},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := New().Classify(request(test.line, contextRoot, grammar.ContentDSL))
			if len(decision.Problems) == 0 || decision.Problems[0].Kind != test.kind {
				t.Fatalf("problems = %#v, want first %q", decision.Problems, test.kind)
			}
		})
	}
}

// TestEveryTagProducesTypedElements проверяет основные элементы репрезентативных форм.
func TestEveryTagProducesTypedElements(t *testing.T) {
	tests := []struct {
		line     string
		lineType model.LineType
	}{
		{line: "@dsl-version 1.2", lineType: model.LineTag},
		{line: "@document-id lesson", lineType: model.LineTag},
		{line: "@resource-dir images, \"audio files\"", lineType: model.LineTag},
		{line: "@media audio \"track 01.mp3\"", lineType: model.LineTag},
		{line: "@example sample", lineType: model.LineTag},
		{line: "@text Title {", lineType: model.LineBlockStart},
		{line: "@multifill", lineType: model.LineTag},
		{line: "@variants {", lineType: model.LineBlockStart},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			decision := New().Classify(request(test.line, contextRoot, grammar.ContentDSL))
			if decision.LineType != test.lineType || len(decision.Elements) == 0 || decision.Elements[0].Type != model.ElementTag {
				t.Fatalf("decision = %#v", decision)
			}
		})
	}
}

// TestResourceDirs проверяет quoted/unquoted пути, запятые, остаток и диапазоны.
func TestResourceDirs(t *testing.T) {
	tests := []struct {
		name, input string
		values      []string
		ranges      []span
		remainder   *span
	}{
		{name: "empty", input: ""},
		{name: "one", input: "images", values: []string{"images"}, ranges: []span{{0, 6}}},
		{name: "spaces", input: " images , audio ", values: []string{"images", "audio"}, ranges: []span{{1, 7}, {10, 15}}},
		{name: "quoted comma", input: `"images, old", audio`, values: []string{"images, old", "audio"}, ranges: []span{{0, 13}, {15, 20}}},
		{name: "literal backslash", input: `assets\images`, values: []string{`assets\images`}, ranges: []span{{0, 13}}},
		{name: "trailing comma", input: "images,", values: []string{"images"}, ranges: []span{{0, 6}}, remainder: &span{6, 7}},
		{name: "unterminated quote", input: `"images`, remainder: &span{0, 7}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := scanResourcePaths(test.input)
			if len(got.paths) != len(test.values) || (got.remainder == nil) != (test.remainder == nil) {
				t.Fatalf("scanResourcePaths(%q) = %#v", test.input, got)
			}
			for index, want := range test.values {
				if got.paths[index].value != want {
					t.Fatalf("path %d = %q, want %q", index, got.paths[index].value, want)
				}
				requireSpan(t, got.paths[index].span, test.ranges[index].start, test.ranges[index].end)
			}
			if test.remainder != nil {
				requireSpan(t, *got.remainder, test.remainder.start, test.remainder.end)
			}
		})
	}

	decision := New().Classify(request(`@resource-dir images, "audio files"`, contextRoot, grammar.ContentDSL))
	if len(decision.Elements) != 3 {
		t.Fatalf("resource-dir elements = %#v", decision.Elements)
	}
	for _, element := range decision.Elements[1:] {
		if element.Type != model.ElementResourcePath {
			t.Fatalf("resource-dir element = %#v", element)
		}
	}
}

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
