package app

import "time"

// Компиляционные проверки фиксируют подменяемые границы оркестратора.
var (
	_ Clock                 = architectureClock{}
	_ ProcessingIDGenerator = architectureIDs{}
	_ FileProcessor         = architectureProcessor{}
)

// architectureClock представляет тестовую границу времени для будущих поведенческих тестов.
type architectureClock struct{}

// Now соответствует контракту Clock.
func (architectureClock) Now() time.Time { panic("TODO") }

// architectureIDs представляет тестовый генератор processing ID.
type architectureIDs struct{}

// NewProcessingID соответствует контракту ProcessingIDGenerator.
func (architectureIDs) NewProcessingID() string { panic("TODO") }

// architectureProcessor представляет тестовую границу обработки файла.
type architectureProcessor struct{}

// Process соответствует контракту FileProcessor.
func (architectureProcessor) Process(request ProcessRequest) FileResult { panic("TODO") }
