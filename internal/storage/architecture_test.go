package storage

// Компиляционные проверки фиксируют единый стандартный адаптер чтения и записи.
var (
	_ Reader = FileSystem{}
	_ Writer = FileSystem{}
	_        = fileSystemWithOperations
)
