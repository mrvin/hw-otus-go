package handler

import (
	"context"
)

func GetUserName(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if userName, ok := ctx.Value("username").(string); ok {
		return userName
	}
	return ""
}
