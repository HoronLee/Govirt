// Package image 模型
package image

import (
	"govirt/app/models"
	"govirt/pkg/database"
)

type Image struct {
	models.BaseModel

	Name        string `gorm:"unique"`
	UUID        string `gorm:"unique"`
	Type        string
	Size        uint64
	OS          string
	Arch        string
	Source      string
	Status      string
	PoolName    string
	VolumeName  string
	Description string
	Checksum    string
	MinDisk     uint64
	MinRam      uint64

	models.CommonTimestampsField
}

func (image *Image) Create() (rowsAffected int64, err error) {
	result := database.DB.Create(&image)
	return result.RowsAffected, result.Error
}

func (image *Image) Save() (rowsAffected int64, err error) {
	result := database.DB.Save(&image)
	return result.RowsAffected, result.Error
}

func (image *Image) Delete() (rowsAffected int64, err error) {
	result := database.DB.Delete(&image)
	return result.RowsAffected, result.Error
}
