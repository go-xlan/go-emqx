package utils

import (
	"encoding/hex"

	"github.com/google/uuid"
)

func NewUUID() string {
	u := uuid.New()
	return hex.EncodeToString(u[:])
}

func GetValuePointer[T any](v T) *T {
	return &v
}

func GetPointerValue[T any](v *T) T {
	if v != nil {
		return *v
	} else {
		var zero T
		return zero
	}
}
