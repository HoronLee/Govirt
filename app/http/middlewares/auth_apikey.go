// Package middlewares Gin 中间件
package middlewares

import (
	"govirt/app/models/apikey"
	"govirt/pkg/database"
	"govirt/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Bearer Token 中提取凭证
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证凭证")
			return
		}

		// 解析 Bearer Token 格式：Bearer {name}:{raw_key}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "凭证格式错误，请使用 Bearer 开头")
			return
		}

		credential := strings.SplitN(parts[1], ":", 2)
		if len(credential) != 2 {
			response.Unauthorized(c, "凭证格式错误，请使用 name:key 格式")
			return
		}
		name, rawKey := credential[0], credential[1]

		// 通过 name 查询数据库
		var apikeyModel apikey.Apikey
		if err := database.DB.Where("name = ?", name).First(&apikeyModel).Error; err != nil {
			response.Unauthorized(c, "API Key 不存在")
			return
		}

		// 验证密钥
		if !apikeyModel.CompareApikey(rawKey) {
			response.Unauthorized(c, "API Key 验证失败")
			return
		}

		c.Next()
	}
}
