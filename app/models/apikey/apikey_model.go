package apikey

import (
	"gohub/app/models"
	"gohub/pkg/database"
	"gohub/pkg/hash"
)

type Apikey struct {
	models.BaseModel
	Name string `json:"name,omitempty"`
	Key  string `json:"-"`
	models.CommonTimestampsField
}

func (apikeyModel *Apikey) Create() {
	database.DB.Create(&apikeyModel)
}

func (apikeyModel *Apikey) Delete() int64 {
	result := database.DB.Delete(&apikeyModel)
	return result.RowsAffected
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
