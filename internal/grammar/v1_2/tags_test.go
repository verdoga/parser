package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"reflect"
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

// TestTagSpecsContainClosedDictionary проверяет полноту, порядок, уникальность и независимость закрытого словаря.
func TestTagSpecsContainClosedDictionary(t *testing.T) {
	want := []tagName{
		tagDSLVersion, tagDocumentID, tagSection, tagOrder, tagResourceDir,
		tagHeader, tagTask, tagEndTask, tagStep, tagSpeaking, tagNewPage,
		tagEditor, tagMedia, tagExample, tagWordlist, tagTable, tagScript,
		tagText, tagKey, tagInstruction, tagNote, tagAlt, tagHint, tagFragment,
		tagInclude, tagAnswer, tagQuestion, tagMultifill, tagChoice,
		tagMultichoice, tagMatching, tagOrdering, tagVariants, tagVariant,
	}
	specs := tagSpecs()
	if len(specs) != supportedTagCount {
		t.Fatalf("tagSpecs() length = %d, want %d", len(specs), supportedTagCount)
	}
	got := make([]tagName, len(specs))
	for i := range specs {
		got[i] = specs[i].name
		if specs[i].forms == 0 {
			t.Errorf("tag %q has no forms", specs[i].name)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tag names = %q, want %q", got, want)
	}
	specs[0].name = "changed"
	if next := tagSpecs(); len(next) == 0 || next[0].name != tagDSLVersion {
		t.Fatal("tagSpecs() exposes shared mutable storage")
	}
}

// TestFindTagSpec проверяет поиск каждого канонического имени и отказ для неизвестных имён.
func TestFindTagSpec(t *testing.T) {
	for _, want := range tagSpecs() {
		got, ok := findTagSpec(want.name)
		if !ok || !reflect.DeepEqual(got, want) {
			t.Errorf("findTagSpec(%q) = (%#v, %t), want %#v", want.name, got, ok, want)
		}
	}
	for _, name := range []tagName{"", "TASK", "unknown"} {
		if got, ok := findTagSpec(name); ok || !reflect.DeepEqual(got, tagSpec{}) {
			t.Errorf("findTagSpec(%q) = (%#v, %t), want zero, false", name, got, ok)
		}
	}
}

// TestScanTag проверяет байтовые границы и канонизацию имени тега.
func TestScanTag(t *testing.T) {
	tests := []struct {
		value string
		span  span
		name  tagName
		ok    bool
	}{
		{value: "@task id", span: span{0, 5}, name: tagTask, ok: true},
		{value: "@DSL-Version\t1.2", span: span{0, 12}, name: tagDSLVersion, ok: true},
		{value: "@resource-dir", span: span{0, 13}, name: tagResourceDir, ok: true},
		{value: "@unknown! tail", span: span{0, 8}, name: "unknown", ok: true},
		{value: "task id"}, {value: "@"}, {value: "@ bad"},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			location, name, ok := scanTag(test.value)
			if location != test.span || name != test.name || ok != test.ok {
				t.Fatalf("scanTag() = (%#v, %q, %t), want (%#v, %q, %t)", location, name, ok, test.span, test.name, test.ok)
			}
		})
	}
}

// TestDetectForm проверяет однострочную, блочную и неоднозначную формы.
func TestDetectForm(t *testing.T) {
	tests := []struct {
		value     string
		form      grammar.TagForm
		blockOpen *span
		ambiguous bool
	}{
		{value: "@note text", form: grammar.TagFormLine},
		{value: "@note text {", form: grammar.TagFormBlock, blockOpen: spanPointer(span{11, 12})},
		{value: "@note text { \t", form: grammar.TagFormBlock, blockOpen: spanPointer(span{11, 12})},
		{value: "@note { tail", ambiguous: true},
		{value: "@note {{", ambiguous: true},
		{value: `@note \{literal\}`, form: grammar.TagFormLine},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got := detectForm(test.value)
			if got.form != test.form || got.ambiguous != test.ambiguous || !reflect.DeepEqual(got.blockOpen, test.blockOpen) {
				t.Fatalf("detectForm(%q) = %#v", test.value, got)
			}
		})
	}
}

// TestElementDecisionConstructors проверяет обычный, отсутствующий, неразобранный и остаточный элементы.
func TestElementDecisionConstructors(t *testing.T) {
	value := "name"
	if got := elementDecision(model.ElementName, span{2, 6}, &value); got.Type != model.ElementName || got.Value == nil || *got.Value != value || got.ByteStart != 2 || got.ByteEnd != 6 {
		t.Fatalf("elementDecision() = %#v", got)
	}
	if got := missingElementDecision(model.ElementIdentifier, 7); got.Type != model.ElementIdentifier || got.Value == nil || *got.Value != "" || got.ByteStart != 7 || got.ByteEnd != 7 {
		t.Fatalf("missingElementDecision() = %#v", got)
	}
	if got := unparsedElementDecision(span{8, 10}); got.Type != model.ElementUnparsed || got.Value != nil || got.ByteStart != 8 || got.ByteEnd != 10 {
		t.Fatalf("unparsedElementDecision() = %#v", got)
	}
	if got := contentRemainderDecision(span{11, 15}, "tail"); got.Type != model.ElementContent || got.Value == nil || *got.Value != "tail" || got.ByteStart != 11 || got.ByteEnd != 15 {
		t.Fatalf("contentRemainderDecision() = %#v", got)
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
