package v1_2

import (
	"dslparser/internal/grammar"
	"dslparser/internal/model"
	"testing"
)

// Компиляционные проверки фиксируют контексты и решения состояния DSL v1.2.
var (
	_ grammar.Context                                                        = contextRoot
	_ grammar.Context                                                        = contextOpaque
	_ grammar.Context                                                        = contextExample
	_ grammar.Context                                                        = contextWordlist
	_ grammar.Context                                                        = contextTable
	_ grammar.Context                                                        = contextText
	_ grammar.Context                                                        = contextEditor
	_ grammar.Context                                                        = contextInstruction
	_ grammar.Context                                                        = contextFragment
	_ grammar.Context                                                        = contextChoice
	_ grammar.Context                                                        = contextMatching
	_ grammar.Context                                                        = contextMultifill
	_ grammar.Context                                                        = contextVariants
	_ func(tagSpec, grammar.TagForm) grammar.Transition                      = transitionFor
	_ func(tagSpec, grammar.TagForm, grammar.Context) grammar.ParentDecision = parentFor
	_ func(tagSpec) grammar.BlockDecision                                    = blockFor
)

// TestHeadings проверяет уровни, названия, ошибки формы, экранирование и корневого родителя.
func TestHeadings(t *testing.T) {
	tests := []struct {
		name, line, level, title string
		valid                    bool
	}{
		{name: "level one", line: "# Unit", valid: true, level: "1", title: "Unit"},
		{name: "level two", line: "## Topic", valid: true, level: "2", title: "Topic"},
		{name: "level three", line: "### Step", valid: true, level: "3", title: "Step"},
		{name: "missing space", line: "#Title"},
		{name: "missing title", line: "## "},
		{name: "extra level", line: "#### Title"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scanned := scanHeading(test.line)
			if scanned.valid != test.valid || scanned.value != test.level {
				t.Fatalf("scanHeading(%q) = %#v", test.line, scanned)
			}
			decision := classifyHeading(request(test.line, contextRoot, grammar.ContentDSL))
			if test.valid {
				if decision.LineType != model.LineHeading || len(decision.Elements) != 2 || len(decision.Parent.Kinds) != 0 {
					t.Fatalf("heading decision = %#v", decision)
				}
				if decision.Elements[0].Value == nil || *decision.Elements[0].Value != test.level || decision.Elements[1].Value == nil || *decision.Elements[1].Value != test.title {
					t.Fatalf("heading elements = %#v", decision.Elements)
				}
			} else if decision.LineType != model.LineInvalid || len(decision.Problems) == 0 {
				t.Fatalf("invalid heading decision = %#v", decision)
			}
		})
	}
	if got := New().Classify(request(`\# literal`, contextRoot, grammar.ContentDSL)); got.LineType != model.LineContent {
		t.Fatalf("escaped heading type = %q, want content", got.LineType)
	}
}

// TestContexts проверяет классификацию blank, content, separator и границ во всех режимах.
func TestContexts(t *testing.T) {
	tests := []struct {
		name    string
		context grammar.Context
		mode    grammar.ContentMode
		line    string
		want    model.LineType
	}{
		{name: "root blank", context: contextRoot, mode: grammar.ContentDSL, line: "", want: model.LineBlank},
		{name: "root content", context: contextRoot, mode: grammar.ContentDSL, line: "ordinary", want: model.LineContent},
		{name: "opaque tag literal", context: contextOpaque, mode: grammar.ContentOpaque, line: "@task id", want: model.LineContent},
		{name: "example separator", context: contextExample, mode: grammar.ContentOpaque, line: "---", want: model.LineSeparator},
		{name: "wordlist content", context: contextWordlist, mode: grammar.ContentOpaque, line: "one; two", want: model.LineContent},
		{name: "table separator is content", context: contextTable, mode: grammar.ContentOpaque, line: "---", want: model.LineContent},
		{name: "text blank", context: contextText, mode: grammar.ContentOpaque, line: "", want: model.LineBlank},
		{name: "editor heading literal", context: contextEditor, mode: grammar.ContentEditor, line: "### note", want: model.LineContent},
		{name: "instruction hint", context: contextInstruction, mode: grammar.ContentDSL, line: "@hint help", want: model.LineTag},
		{name: "fragment content", context: contextFragment, mode: grammar.ContentDSL, line: "body", want: model.LineContent},
		{name: "choice answer", context: contextChoice, mode: grammar.ContentDSL, line: "@answer yes", want: model.LineTag},
		{name: "matching separator", context: contextMatching, mode: grammar.ContentDSL, line: "---", want: model.LineSeparator},
		{name: "multifill blank content", context: contextMultifill, mode: grammar.ContentMultifill, line: "", want: model.LineContent},
		{name: "variants branch", context: contextVariants, mode: grammar.ContentDSL, line: "@variant A", want: model.LineTag},
		{name: "block close", context: contextOpaque, mode: grammar.ContentOpaque, line: "}", want: model.LineBlockEnd},
		{name: "malformed close", context: contextRoot, mode: grammar.ContentDSL, line: "} tail", want: model.LineInvalid},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decision := New().Classify(request(test.line, test.context, test.mode))
			if decision.LineType != test.want {
				t.Fatalf("Classify(%q) type = %q, want %q; decision=%#v", test.line, decision.LineType, test.want, decision)
			}
		})
	}
}

// TestTransitionsParentsAndBlocks проверяет изменения областей и открытие блока.
func TestTransitionsParentsAndBlocks(t *testing.T) {
	tests := []struct {
		line  string
		check func(grammar.Transition) bool
	}{
		{line: "@task id", check: func(got grammar.Transition) bool { return got.OpenTask }},
		{line: "@endtask", check: func(got grammar.Transition) bool { return got.CloseTask }},
		{line: "@step Part", check: func(got grammar.Transition) bool { return got.SetInnerStep }},
		{line: "@variants {", check: func(got grammar.Transition) bool { return got.OpenVariants && got.CloseTask }},
		{line: "@variant A", check: func(got grammar.Transition) bool { return got.SetVariant }},
	}
	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			got := New().Classify(request(test.line, contextRoot, grammar.ContentDSL))
			if !test.check(got.Transition) {
				t.Fatalf("transition = %#v", got.Transition)
			}
		})
	}
	block := New().Classify(request("@text {", contextRoot, grammar.ContentDSL)).Block
	if block.Action != grammar.BlockOpen || block.Tag != "text" || block.Context != contextText || block.ContentMode != grammar.ContentOpaque {
		t.Fatalf("text block = %#v", block)
	}
}
