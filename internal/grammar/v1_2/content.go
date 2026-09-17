package v1_2

// unescape снимает только экранирование, значимое для заданной политики содержимого.
func unescape(value string, policy contentPolicy) string { panic("TODO") }

// unescapeText снимает только экранирование, имеющее специальное значение в обычном тексте.
func unescapeText(value string) string { panic("TODO") }

// unescapeMediaSource снимает экранирование кавычки и обратной косой черты источника media.
func unescapeMediaSource(value string) string { panic("TODO") }

// unescapedBraceRanges возвращает диапазоны структурно незащищённых фигурных скобок.
func unescapedBraceRanges(value string, protected []span) []span { panic("TODO") }
