// filepath: /home/horonlee/projects/govirt/pkg/libvirtd/image.go
package libvirtd

import (
	"errors"
	"fmt"
	imageMod "govirt/app/models/image"
	"govirt/pkg/database"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"
	"os"
	"strings"
)

// 定义镜像状态常量
const (
	ImageStatusCreating = "creating"
	ImageStatusActive   = "active"
	ImageStatusError    = "error"
	ImageStatusDeleting = "deleting"
	ImageStatusDeleted  = "deleted"
)

// CreateImageFromLocalFile 从本地文件创建镜像
func (vc *VirtConn) CreateImageFromLocalFile(name, sourceFilePath, poolName, osType, arch, imageType, description string, minDisk, minRam uint64) (*imageMod.Image, error) {
	// 1. 基本检查
	if name == "" || sourceFilePath == "" || poolName == "" {
		return nil, fmt.Errorf("名称、源文件路径和存储池名称不能为空")
	}

	// 2. 检查源文件是否存在
	if _, err := os.Stat(sourceFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("源文件不存在: %v", err)
	}

	// 3. 获取文件大小
	fileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return nil, fmt.Errorf("获取源文件信息失败: %v", err)
	}
	fileSize := uint64(fileInfo.Size())

	// 4. 检查是否已存在同名镜像
	if imageMod.IsExist("name", name) {
		return nil, fmt.Errorf("已存在同名镜像: %s", name)
	}

	// 5. 获取存储池
	pool, err := vc.GetStoragePool(poolName)
	if err != nil {
		return nil, fmt.Errorf("获取存储池失败: %v", err)
	}

	// 6. 生成UUID和Volume名称
	uuid := helpers.GenerateUUIDString()
	volumeName := fmt.Sprintf("%s_%s.%s", name, uuid, imageType)

	// 7. 创建Image记录
	image := &imageMod.Image{
		Name:        name,
		UUID:        uuid,
		Type:        imageType,
		Size:        fileSize,
		OS:          osType,
		Arch:        arch,
		Source:      sourceFilePath,
		Status:      ImageStatusCreating,
		PoolName:    poolName,
		VolumeName:  volumeName,
		Description: description,
		MinDisk:     minDisk,
		MinRam:      minRam,
	}

	// 8. 保存到数据库
	if _, err := image.Create(); err != nil {
		return nil, fmt.Errorf("创建镜像记录失败: %v", err)
	}

	// 9. 创建存储卷
	volParams := &xmlDefine.VolumeTemplateParams{
		Name:     volumeName,
		Capacity: minDisk,
		Type:     imageType,
	}

	vol, err := vc.CreateVolume(pool, volParams, 0)
	if err != nil {
		// 创建卷失败，更新状态并返回错误
		image.Status = ImageStatusError
		image.Save()
		return nil, fmt.Errorf("创建存储卷失败: %v", err)
	}

	// 10. 上传文件到存储卷
	file, err := os.Open(sourceFilePath)
	if err != nil {
		// 打开文件失败，删除卷并更新状态
		vc.Libvirt.StorageVolDelete(vol, 0)
		image.Status = ImageStatusError
		image.Save()
		return nil, fmt.Errorf("打开源文件失败: %v", err)
	}
	defer file.Close()

	// 上传文件到存储卷
	err = vc.Libvirt.StorageVolUpload(vol, file, 0, fileSize, 0)
	if err != nil {
		// 上传失败，删除卷并更新状态
		vc.Libvirt.StorageVolDelete(vol, 0)
		image.Status = ImageStatusError
		image.Save()
		return nil, fmt.Errorf("上传文件到存储卷失败: %v", err)
	}

	// 11. 生成校验和
	checksum, err := helpers.CalculateChecksum(sourceFilePath)
	if err != nil {
		// 校验和计算失败，记录错误但继续
		fmt.Printf("计算校验和失败，但会继续: %v\n", err)
	} else {
		image.Checksum = checksum
		image.Save()
	}

	// 12. 更新状态为活动状态
	image.Status = ImageStatusActive
	if _, err := image.Save(); err != nil {
		return nil, fmt.Errorf("更新镜像状态失败: %v", err)
	}

	return image, nil
}

// CreateImageFromURL 从URL创建镜像
func (vc *VirtConn) CreateImageFromURL(name, sourceURL, poolName, osType, arch, imageType, description string, minDisk, minRam uint64) (*imageMod.Image, error) {
	// 这里可以实现从URL下载文件到临时目录，然后调用CreateImageFromLocalFile，最后清理临时文件
	// 由于涉及到网络下载，可能需要支持断点续传、进度显示等功能
	// 为了简化示例，这里省略具体实现
	return nil, fmt.Errorf("从URL创建镜像功能尚未实现")
}

// DeleteImage 删除镜像及关联的存储卷
func (vc *VirtConn) DeleteImage(idOrUUID string) error {
	var image imageMod.Image
	var err error

	// 1. 查找镜像记录
	if helpers.IsUUIDString(idOrUUID) {
		image, err = imageMod.GetByUUID(idOrUUID)
	} else {
		image, err = imageMod.Get(idOrUUID)
	}

	if err != nil {
		return fmt.Errorf("查找镜像失败: %v", err)
	}

	// 2. 更新镜像状态为删除中
	image.Status = ImageStatusDeleting
	if _, err := image.Save(); err != nil {
		return fmt.Errorf("更新镜像状态失败: %v", err)
	}

	// 3. 获取存储池
	pool, err := vc.GetStoragePool(image.PoolName)
	if err != nil {
		image.Status = ImageStatusError
		image.Save()
		return fmt.Errorf("获取存储池失败: %v", err)
	}

	// 4. 删除存储卷
	err = vc.DeleteVolume(pool, image.VolumeName, 0)
	if err != nil {
		// 如果错误是"no storage vol with matching"，表示卷可能已经被删除
		if strings.Contains(err.Error(), "no storage vol with matching") {
			logger.WarnString("image", "删除镜像", fmt.Sprintf("卷 %s 不存在，可能已被删除", image.VolumeName))
		} else {
			image.Status = ImageStatusError
			image.Save()
			return fmt.Errorf("删除存储卷失败: %v", err)
		}
	}

	// 5. 删除数据库记录
	if _, err := image.Delete(); err != nil {
		return fmt.Errorf("删除镜像记录失败: %v", err)
	}

	return nil
}

// GetImagePath 获取镜像对应的存储卷路径
func (vc *VirtConn) GetImagePath(image *imageMod.Image) (string, error) {
	pool, err := vc.GetStoragePool(image.PoolName)
	if err != nil {
		return "", fmt.Errorf("获取存储池失败: %v", err)
	}

	vol, err := vc.GetVolume(pool, image.VolumeName)
	if err != nil {
		return "", fmt.Errorf("获取卷失败: %v", err)
	}

	path, err := vc.Libvirt.StorageVolGetPath(vol)
	if err != nil {
		return "", fmt.Errorf("获取卷路径失败: %v", err)
	}

	return path, nil
}

// ListActiveImages 列出所有活动状态的镜像
// flag: 0表示返回summary（只包含Name、UUID、Type），1表示返回details（完整信息）
func (vc *VirtConn) ListActiveImages(flag int) (any, error) {
	if flag == 0 {
		var summaries []struct {
			Name string
			UUID string
			Type string
		}
		result := database.DB.Model(&imageMod.Image{}).Select("name, uuid, type").Where("status = ?", ImageStatusActive).Find(&summaries)
		return summaries, result.Error
	} else if flag == 1 {
		var images []*imageMod.Image
		result := database.DB.Where("status = ?", ImageStatusActive).Find(&images)
		return images, result.Error
	}
	return nil, fmt.Errorf("无效的flag值: %d", flag)
}

// SyncImagesWithVolumes 同步数据库中的镜像记录与实际存储卷
// 如果数据库中有记录但存储卷不存在，则将镜像状态更新为deleted
// 如果数据库中的记录状态为deleted但卷已经存在，则恢复状态为active
// 不删除任何池中的卷
func (vc *VirtConn) SyncImagesWithVolumes(poolName string) error {
	// 1. 获取存储池
	pool, err := vc.GetStoragePool(poolName)
	if err != nil {
		return fmt.Errorf("获取存储池失败: %v", err)
	}

	// 2. 获取池中所有卷的列表
	volumeNames, err := vc.ListVolumesSummary(pool)
	if err != nil {
		return fmt.Errorf("获取卷列表失败: %v", err)
	}

	// 创建存储卷的映射表，用于快速查询
	volumeMap := make(map[string]bool)
	for _, name := range volumeNames {
		volumeMap[name] = true
	}

	// 3. 获取数据库中的所有镜像记录
	dbImages, err := imageMod.All()
	if err != nil {
		return fmt.Errorf("获取镜像记录失败: %v", err)
	}

	// 4. 遍历所有数据库记录，检查其关联的存储卷是否存在
	for i := range dbImages {
		// 仅处理与当前存储池匹配的镜像记录
		if dbImages[i].PoolName != poolName {
			continue
		}

		volumeExists := volumeMap[dbImages[i].VolumeName]

		if dbImages[i].Status != ImageStatusDeleted && !volumeExists {
			// 镜像记录状态不是deleted但卷不存在，更新状态为deleted
			dbImages[i].Status = ImageStatusDeleted
			if _, err := dbImages[i].Save(); err != nil {
				logger.ErrorString("image", "更新镜像信息", fmt.Sprintf("更新镜像 %s 状态失败: %v\n", dbImages[i].Name, err))
				continue
			}
			logger.WarnString("image", "更新镜像信息", fmt.Sprintf("镜像 %s 的存储卷 %s 不存在，已将状态更新为deleted\n",
				dbImages[i].Name, dbImages[i].VolumeName))
		} else if dbImages[i].Status == ImageStatusDeleted && volumeExists {
			// 镜像记录状态是deleted但卷已存在，恢复状态为active
			dbImages[i].Status = ImageStatusActive
			if _, err := dbImages[i].Save(); err != nil {
				logger.ErrorString("image", "更新镜像信息", fmt.Sprintf("更新镜像 %s 状态失败: %v\n", dbImages[i].Name, err))
				continue
			}
			logger.InfoString("image", "更新镜像信息", fmt.Sprintf("镜像 %s 的存储卷 %s 已恢复，状态更新为active\n",
				dbImages[i].Name, dbImages[i].VolumeName))
		}
	}

	return nil
}

// SyncAllImagesWithVolumes 同步所有存储池中的镜像记录与实际存储卷
func (vc *VirtConn) SyncAllImagesWithVolumes() error {
	// 获取数据库中所有镜像所在的存储池列表（去重）
	var poolNames []string
	if err := database.DB.Model(&imageMod.Image{}).Distinct("pool_name").Pluck("pool_name", &poolNames).Error; err != nil {
		return fmt.Errorf("获取存储池列表失败: %v", err)
	}

	// 为每个池执行同步操作
	var syncErrors []error
	for _, poolName := range poolNames {
		if err := vc.SyncImagesWithVolumes(poolName); err != nil {
			syncErrors = append(syncErrors, fmt.Errorf("同步池 %s 失败: %v", poolName, err))
		}
	}

	// 如果有错误，返回所有错误信息
	if len(syncErrors) > 0 {
		errMsg := "同步过程中发生以下错误:\n"
		for _, err := range syncErrors {
			errMsg += "- " + err.Error() + "\n"
		}
		return errors.New(errMsg)
	}

	return nil
}
