package parser

import "dslparser/internal/model"

// canRecover сообщает true при однозначном продолжении и false при необходимости остановки.
func canRecover(decision RecoveryDecision) bool { panic("TODO") }

// appendLexicalTail добавляет все строки начиная с from в безопасном лексическом режиме.
func appendLexicalTail(lines []model.Line, source []physicalLine, from int) []model.Line {
	panic("TODO")
}

// lexicalTailLine строит корневую blank- или content-строку без снятия экранирования.
func lexicalTailLine(line physicalLine) model.Line { panic("TODO") }

// unclosedBlockDiagnostic создаёт фатальную P011 со связанным местом открытия блока.
func unclosedBlockDiagnostic(frame blockFrame, eofLine int, builder *diagnosticBuilder) model.Diagnostic {
	panic("TODO")
}
