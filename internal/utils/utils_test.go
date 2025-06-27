package utils_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-xlan/go-emqx/internal/utils"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
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

func TestGetAccountsTokens(t *testing.T) {
	accounts := gin.Accounts{
		"username-abc": "password-123",
	}
	t.Log(neatjsons.S(utils.GetAccountsTokens(accounts)))
}
