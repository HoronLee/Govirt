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

func (apikeyModel *Apikey) Create() {
	database.DB.Create(&apikeyModel)
}

func (apikeyModel *Apikey) Delete() int64 {
	result := database.DB.Delete(&apikeyModel)
	return result.RowsAffected
}
