package image

import (
	"govirt/pkg/database"
)

func Get(idstr string) (image Image, err error) {
	result := database.DB.Where("id", idstr).First(&image)
	return image, result.Error
}

func GetBy(field, value string) (image Image, err error) {
	result := database.DB.Where("? = ?", field, value).First(&image)
	return image, result.Error
}

func GetByUUID(uuid string) (image Image, err error) {
	result := database.DB.Where("uuid = ?", uuid).First(&image)
	return image, result.Error
}

func All() (images []Image, err error) {
	result := database.DB.Find(&images)
	return images, result.Error
}

func IsExist(field, value string) bool {
	var count int64
	database.DB.Model(Image{}).Where("? = ?", field, value).Count(&count)
	return count > 0
}
