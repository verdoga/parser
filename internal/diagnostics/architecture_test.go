package diagnostics

import (
	"errors"

	"dslparser/internal/model"
)

// Компиляционные проверки фиксируют внешний контракт реестра и накопителя диагностик.
var (
	_ error = (*RequestError)(nil)
	_       = Request{
		Code:             P011,
		Scope:            model.ScopeBlock,
		Fatal:            true,
		Location:         &model.Range{},
		RelatedLocations: []model.Range{},
		Details:          Details{Tag: "text"},
	}
	_ = Details{Received: "2.0", Supported: []string{"1.2"}, Path: "/lesson.txt", Cause: errors.New("permission denied")}
)
