package grammar

// RegistrationErrorReason задаёт устойчивую причину невозможности построить реестр.
type RegistrationErrorReason int

const (
	// RegistrationNilGrammar означает, что вместо реализации передано nil.
	RegistrationNilGrammar RegistrationErrorReason = iota
	// RegistrationEmptyVersion означает, что грамматика вернула пустую версию.
	RegistrationEmptyVersion
	// RegistrationDuplicateVersion означает, что версия уже зарегистрирована.
	RegistrationDuplicateVersion
)

// RegistrationError описывает недопустимую запись при построении реестра.
type RegistrationError struct {
	// version содержит значение версии, связанное с ошибкой.
	version string
	// reason содержит устойчивую причину ошибки.
	reason RegistrationErrorReason
}

// Error возвращает человекочитаемое описание ошибки регистрации.
func (e *RegistrationError) Error() string { panic("TODO") }

// Version возвращает версию, связанную с ошибкой.
func (e *RegistrationError) Version() string { panic("TODO") }

// Reason возвращает устойчивую причину ошибки регистрации.
func (e *RegistrationError) Reason() RegistrationErrorReason { panic("TODO") }
