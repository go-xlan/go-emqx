package utils

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yyle88/must"
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

func GetAccountsTokens(accounts gin.Accounts) map[string]string {
	res := make(map[string]string, len(accounts))
	for acc, pwd := range accounts {
		authToken := must.Nice(acc) + ":" + must.Nice(pwd)
		token := "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken))
		res[acc] = token
	}
	return res
}
