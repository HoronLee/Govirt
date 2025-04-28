package image

import (
	"fmt"
	imageMod "govirt/app/models/image"
	"govirt/pkg/database"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
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
// flag: 0表示返回summary（只包含Name、UUID、Type），1表示返回details（完整信息）
func ListActiveImages(flag int) (any, error) {
	if flag == 0 {
		var summaries []struct {
			Name string
			UUID string
			Type string
		}
		result := database.DB.Model(&imageMod.Image{}).Select("name, uuid, type").Where("status = ?", StatusActive).Find(&summaries)
		return summaries, result.Error
	} else if flag == 1 {
		var images []*imageMod.Image
		result := database.DB.Where("status = ?", StatusActive).Find(&images)
		return images, result.Error
	}
	return nil, fmt.Errorf("无效的flag值: %d", flag)
}

// SyncImagesWithVolumes 同步数据库中的镜像记录与实际存储卷
// 如果数据库中有记录但存储卷不存在，则将镜像状态更新为deleted
// 不删除任何池中的卷
func SyncImagesWithVolumes(poolName string) error {
	// 1. 获取存储池
	pool, err := storagePool.GetStoragePool(poolName)
	if err != nil {
		return fmt.Errorf("获取存储池失败: %v", err)
	}

	// 2. 从数据库获取所有镜像
	var dbImages []imageMod.Image
	if err := database.DB.Where("pool_name = ?", poolName).Find(&dbImages).Error; err != nil {
		return fmt.Errorf("查询数据库镜像记录失败: %v", err)
	}

	// 3. 获取存储池中所有卷
	volumeNames, err := volume.ListVolumesSummary(pool)
	if err != nil {
		return fmt.Errorf("获取存储池中的卷列表失败: %v", err)
	}

	// 4. 创建卷名到镜像的映射关系
	volumeMap := make(map[string]bool)
	for _, name := range volumeNames {
		volumeMap[name] = true
	}

	// 5. 检查数据库中的镜像记录是否有对应的实际卷
	// 如果没有，则更新状态为deleted
	for i := range dbImages {
		if dbImages[i].Status != StatusDeleted && !volumeMap[dbImages[i].VolumeName] {
			// 镜像记录存在但卷不存在，更新状态为deleted
			dbImages[i].Status = StatusDeleted
			if _, err := dbImages[i].Save(); err != nil {
				logger.InfoString("image", "更新镜像信息", fmt.Sprintf("更新镜像 %s 状态失败: %v\n", dbImages[i].Name, err))
				continue
			}
			logger.WarnString("image", "更新镜像信息", fmt.Sprintf("镜像 %s 的存储卷 %s 不存在，已将状态更新为deleted\n",
				dbImages[i].Name, dbImages[i].VolumeName))
		}
	}

	return nil
}

// SyncAllImagesWithVolumes 同步所有存储池中的镜像记录与实际存储卷
func SyncAllImagesWithVolumes() error {
	// 获取数据库中所有镜像所在的存储池列表（去重）
	var poolNames []string
	if err := database.DB.Model(&imageMod.Image{}).Distinct("pool_name").Pluck("pool_name", &poolNames).Error; err != nil {
		return fmt.Errorf("获取存储池列表失败: %v", err)
	}

	// 为每个池执行同步操作
	var syncErrors []error
	for _, poolName := range poolNames {
		if err := SyncImagesWithVolumes(poolName); err != nil {
			syncErrors = append(syncErrors, fmt.Errorf("同步池 %s 失败: %v", poolName, err))
		}
	}

	// 如果有错误，返回所有错误信息
	if len(syncErrors) > 0 {
		errMsg := "同步过程中发生以下错误:\n"
		for _, err := range syncErrors {
			errMsg += "- " + err.Error() + "\n"
		}
		return fmt.Errorf(errMsg)
	}

	return nil
}
