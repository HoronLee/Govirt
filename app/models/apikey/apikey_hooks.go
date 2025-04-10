package apikey

import (
	"gohub/pkg/hash"

	"gorm.io/gorm"
)

// BeforeSave GORM 的模型钩子，在创建和更新模型前调用
func (apikeyModel *Apikey) BeforeSave(tx *gorm.DB) (err error) {

	if !hash.BcryptIsHashed(apikeyModel.Key) {
		apikeyModel.Key = hash.BcryptHash(apikeyModel.Key)
	}
	return
}
