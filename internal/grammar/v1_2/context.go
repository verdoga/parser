package v1_2

import "dslparser/internal/grammar"

const (
	// contextRoot распознаёт обычные конструкции на уровне документа, шага или задания.
	contextRoot grammar.Context = "v1.2/root"
	// contextOpaque распознаёт непрозрачное содержимое контентного блока.
	contextOpaque grammar.Context = "v1.2/opaque"
	// contextExample распознаёт непрозрачное содержимое example и строку разделителя.
	contextExample grammar.Context = "v1.2/example"
	// contextWordlist распознаёт непрозрачные строки блочного wordlist.
	contextWordlist grammar.Context = "v1.2/wordlist"
	// contextTable сохраняет строки table как содержимое, включая строку из трёх дефисов.
	contextTable grammar.Context = "v1.2/table"
	// contextText сохраняет строки text без структурного разбора HTML.
	contextText grammar.Context = "v1.2/text"
	// contextEditor распознаёт литеральное содержимое редакторского блока.
	contextEditor grammar.Context = "v1.2/editor"
	// contextInstruction распознаёт содержимое блочного instr и вложенный hint.
	contextInstruction grammar.Context = "v1.2/instruction"
	// contextFragment распознаёт разрешённое содержимое fragment.
	contextFragment grammar.Context = "v1.2/fragment"
	// contextChoice распознаёт варианты выбора, answer и связанный hint.
	contextChoice grammar.Context = "v1.2/choice"
	// contextMatching распознаёт две части matching или ordering и их разделитель.
	contextMatching grammar.Context = "v1.2/matching"
	// contextMultifill распознаёт предзаполненное непрозрачное содержимое multifill.
	contextMultifill grammar.Context = "v1.2/multifill"
	// contextVariants распознаёт ветви внешнего блока variants.
	contextVariants grammar.Context = "v1.2/variants"
)

// classifyBlank создаёт решение для пустой после обрезки строки.
func classifyBlank(request grammar.GrammarRequest) grammar.GrammarDecision { panic("TODO") }

// classifyBlockBoundary создаёт решение для отдельной или неправильной границы блока.
func classifyBlockBoundary(request grammar.GrammarRequest) grammar.GrammarDecision {
	panic("TODO")
}

// classifyDeclaration создаёт решение для строки, начинающей распознавание DSL-тега.
func classifyDeclaration(request grammar.GrammarRequest) grammar.GrammarDecision { panic("TODO") }

// classifyContent создаёт решение для обычного или непрозрачного содержимого.
func classifyContent(request grammar.GrammarRequest) grammar.GrammarDecision { panic("TODO") }

// classifySeparator создаёт решение для разрешённой отдельной строки разделителя.
func classifySeparator(request grammar.GrammarRequest) grammar.GrammarDecision { panic("TODO") }
