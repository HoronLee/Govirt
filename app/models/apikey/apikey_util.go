package apikey

import (
	"govirt/pkg/database"
	"govirt/pkg/hash"
)

func GetFromID(idstr string) (apikeyModel Apikey, err error) {
	result := database.DB.Where("id", idstr).First(&apikeyModel)
	return apikeyModel, result.Error
}

func GetFromName(namestr string) (apikeyModel Apikey, err error) {
	result := database.DB.Where("name", namestr).First(&apikeyModel)
	return apikeyModel, result.Error
}
func (apikeyModel *Apikey) CompareApikey(_key string) bool {
	return hash.BcryptCheck(_key, apikeyModel.Key)
}

// IsExist 检查 API Key 是否存在
func (apikey *Apikey) IsExist() bool {
	var count int64
	database.DB.Model(&Apikey{}).Where("name = ?", apikey.Name).Count(&count)
	return count > 0
}
