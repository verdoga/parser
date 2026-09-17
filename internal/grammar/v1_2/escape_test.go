package v1_2

// Компиляционные проверки фиксируют специальные лексические границы DSL v1.2.
var (
	_ func(string) string                = unescapeText
	_ func(string) string                = unescapeMediaSource
	_ func(string, contentPolicy) string = unescape
	_ func(string) resourcePathResult    = scanResourcePaths
	_ func(string) []span                = placeholderRanges
	_ func(string, []span) []span        = unescapedBraceRanges
)
