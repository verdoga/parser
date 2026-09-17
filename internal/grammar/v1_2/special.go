package v1_2

import "dslparser/internal/grammar"

// resourcePathResult содержит синтаксически выделенные пути и остаток списка.
type resourcePathResult struct {
	// paths содержит пути в исходном порядке.
	paths []argument
	// remainder содержит ошибочный остаток либо nil.
	remainder *span
}

// parseMedia разбирает обязательные TYPE и SOURCE объявления media.
func parseMedia(request grammar.GrammarRequest, spec tagSpec) grammar.GrammarDecision {
	panic("TODO")
}

// parseResourceDirs разбирает список буквальных путей объявления resource-dir.
func parseResourceDirs(request grammar.GrammarRequest, spec tagSpec) grammar.GrammarDecision {
	panic("TODO")
}

// parseMultifill разбирает однострочную или блочную форму с необязательной инструкцией.
func parseMultifill(request grammar.GrammarRequest, spec tagSpec) grammar.GrammarDecision {
	panic("TODO")
}

// scanResourcePaths выделяет заключённые и не заключённые в кавычки пути в исходном порядке.
func scanResourcePaths(value string) resourcePathResult { panic("TODO") }

// placeholderRanges возвращает защищённые диапазоны минимально распознанных multifill-плейсхолдеров.
func placeholderRanges(value string) []span { panic("TODO") }
