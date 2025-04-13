package apikey

import (
	"gohub/pkg/database"
)

func GetFromID(idstr string) (apikeyModel Apikey) {
	database.DB.Where("id", idstr).First(&apikeyModel)
	return
}

func GetFromName(namestr string) (apikeyModel Apikey) {
	database.DB.Where("name", namestr).First(&apikeyModel)
	return
}
