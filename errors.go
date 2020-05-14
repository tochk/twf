package twf

import "errors"

var (
	errorNotSlice        = errors.New("items must be slice")
	ErrFksIndexNotExists = errors.New("fks index not exists")
	ErrFksMustBeSLice    = errors.New("fks must me a slice")
	errorNotStruct       = errors.New("item must be struct")
)
