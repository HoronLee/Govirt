package apikey

import (
	"govirt/app/models"
	"govirt/pkg/database"
)

type Apikey struct {
	models.BaseModel
	Name string `json:"name,omitempty"`
	Key  string `json:"-"`
	models.CommonTimestampsField
}

func (apikeyModel *Apikey) Create() error {
	result := database.DB.Create(&apikeyModel)
	return result.Error
}

func (apikeyModel *Apikey) Delete() (int64, error) {
	result := database.DB.Delete(&apikeyModel)
	return result.RowsAffected, result.Error
}
