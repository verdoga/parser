package app

import "time"

// Компиляционные проверки фиксируют подменяемые границы оркестратора.
var (
	_ clockFunc                                      = architectureNow
	_ processingIDFunc                               = architectureID
	_ processFileFunc                                = architectureProcess
	_                                                = processorDependencies{}
	_ processFileFunc                                = processor{}.process
	_ func(Options) (RunResult, error)               = Run
	_ func() (dependencies, error)                   = defaultDependencies
	_ func(Options, dependencies) (RunResult, error) = runWithDependencies
)

// architectureNow представляет тестовую границу времени для будущих поведенческих тестов.
func architectureNow() time.Time { panic("TODO") }

// architectureID представляет тестовый генератор processing ID.
func architectureID() string { panic("TODO") }

// architectureProcess представляет тестовую границу обработки файла.
func architectureProcess(request processRequest) FileResult { panic("TODO") }
