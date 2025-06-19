package utils_test

import (
	"testing"

	"github.com/go-xlan/go-emqx/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestNewUUID(t *testing.T) {
	t.Log(utils.NewUUID())
}

func TestGetValuePointer(t *testing.T) {
	address := utils.GetValuePointer("abc")
	value := utils.GetPointerValue(address)
	require.Equal(t, "abc", value)
}

func TestGetPointerValue(t *testing.T) {
	address := utils.GetValuePointer(200)
	value := utils.GetPointerValue(address)
	require.Equal(t, 200, value)
}
