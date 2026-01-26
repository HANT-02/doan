package utils

import (
	"doan/pkg/constants"
	"fmt"
)

func GenCacheKeyMarkResendVerifyAccount(userId, accountId string) string {
	return fmt.Sprintf("%s:%s:%s", constants.RedisPrefixMarkResendVerifyAccount, userId, accountId)
}
func GenerateUpdatePasswordEmailCacheKey(accountId string) string {
	return fmt.Sprintf("%s:%s", constants.RedisPrefixEmailUpdatePassword, accountId)
}
