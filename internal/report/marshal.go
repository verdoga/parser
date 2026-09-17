package report

import "dslparser/internal/model"

// Marshal проверяет полный результат и сериализует его с отступом в два пробела,
// без HTML-экранирования и с одним завершающим LF.
// Marshal возвращает ошибку с сохранённой причиной, если нарушен инвариант модели
// или JSON не может быть сформирован.
func Marshal(result model.Result) ([]byte, error) { panic("TODO") }
