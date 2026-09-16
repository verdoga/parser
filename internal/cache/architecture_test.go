package cache

// Компиляционная проверка фиксирует стандартный адаптер минимального чтения JSON.
var _ HeaderReader = JSONReader{}
