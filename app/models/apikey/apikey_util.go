package apikey

import (
	"fmt"
	"gohub/pkg/database"
	"gohub/pkg/helpers"
	"gohub/pkg/logger"
)

func GetFromID(idstr string) (apikeyModel Apikey) {
	database.DB.Where("id", idstr).First(&apikeyModel)
	return
}

func GetFromName(namestr string) (apikeyModel Apikey) {
	database.DB.Where("name", namestr).First(&apikeyModel)
	return
}

// 新增初始化函数
func InitApikey() {
	var count int64
	database.DB.Model(&Apikey{}).Count(&count)

	if count == 0 {
		key := helpers.RandomString(64)

		logger.InfoString("Apikey", "生成初始 APIKey", fmt.Sprintf("APIKey: %s", key))

		apikey := Apikey{
			Name: "InitKey",
			Key:  key,
		}
		database.DB.Create(&apikey)
	}
}
