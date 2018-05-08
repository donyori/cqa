package wrapper

import (
	"errors"

	"github.com/donyori/cqa/data/db/generic"
)

type wrappedAccessor struct {
	accessor generic.Accessor
}

var ErrResultTypeWrong error = errors.New("result type is wrong")

func (wa *wrappedAccessor) GetAccessor() generic.Accessor {
	if wa == nil {
		return nil
	}
	return wa.accessor
}
