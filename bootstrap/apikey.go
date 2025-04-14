package bootstrap

import (
	"fmt"
	"govirt/app/models/apikey"
	"govirt/pkg/database"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
)

// 新增初始化函数
func InitApikey() {
	var count int64
	database.DB.Model(&apikey.Apikey{}).Count(&count)

	if count == 0 {
		key := helpers.RandomString(64)

		logger.InfoString("Apikey", "生成初始 APIKey", fmt.Sprintf("APIKey: %s", key))

		apikey := apikey.Apikey{
			Name: "InitKey",
			Key:  key,
		}
		database.DB.Create(&apikey)
	}
}
