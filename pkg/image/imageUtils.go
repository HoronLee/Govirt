package image

import (
	"fmt"
	imageMod "govirt/app/models/image"
	"govirt/pkg/database"
	"govirt/pkg/libvirtd"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
	"image"
)

// GetImagePath 获取镜像对应的存储卷路径
func GetImagePath(image *imageMod.Image) (string, error) {
	pool, err := storagePool.GetStoragePool(image.PoolName)
	if err != nil {
		return "", fmt.Errorf("获取存储池失败: %v", err)
	}

	vol, err := volume.GetVolume(pool, image.VolumeName)
	if err != nil {
		return "", fmt.Errorf("获取卷失败: %v", err)
	}

	path, err := libvirtd.Connection.StorageVolGetPath(vol)
	if err != nil {
		return "", fmt.Errorf("获取卷路径失败: %v", err)
	}

	return path, nil
}

// ListActiveImages 列出所有活动状态的镜像
func ListActiveImages() ([]image.Image, error) {
	var images []image.Image
	result := database.DB.Where("status = ?", StatusActive).Find(&images)
	return images, result.Error
}
