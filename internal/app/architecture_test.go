package app

import "time"

// Компиляционные проверки фиксируют подменяемые границы оркестратора.
var (
	_ ClockFunc                                      = architectureNow
	_ ProcessingIDFunc                               = architectureID
	_ ProcessFileFunc                                = architectureProcess
	_                                                = ProcessorDependencies{}
	_ ProcessFileFunc                                = Processor{}.Process
	_ func() (Dependencies, error)                   = DefaultDependencies
	_ func(Options, Dependencies) (RunResult, error) = RunWithDependencies
)

// architectureNow представляет тестовую границу времени для будущих поведенческих тестов.
func architectureNow() time.Time { panic("TODO") }

// architectureID представляет тестовый генератор processing ID.
func architectureID() string { panic("TODO") }

// architectureProcess представляет тестовую границу обработки файла.
func architectureProcess(request ProcessRequest) FileResult { panic("TODO") }
