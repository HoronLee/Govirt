package image

import (
	"fmt"
	imageMod "govirt/app/models/image"
	"govirt/pkg/helpers"
	"govirt/pkg/libvirtd"
	"govirt/pkg/storagePool"
	"govirt/pkg/volume"
	"govirt/pkg/xmlDefine"
	"os"
)

// 定义常量
const (
	StatusCreating = "creating"
	StatusActive   = "active"
	StatusError    = "error"
	StatusDeleting = "deleting"
	StatusDeleted  = "deleted"
)

// CreateImageFromLocalFile 从本地文件创建镜像
func CreateImageFromLocalFile(name, sourceFilePath, poolName, osType, arch, imageType, description string, minDisk, minRam uint64) (*imageMod.Image, error) {
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
	pool, err := storagePool.GetStoragePool(poolName)
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
		Status:      StatusCreating,
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

	vol, err := volume.CreateVolume(pool, volParams, 0)
	if err != nil {
		// 创建卷失败，更新状态并返回错误
		image.Status = StatusError
		image.Save()
		return nil, fmt.Errorf("创建存储卷失败: %v", err)
	}

	// 10. 上传文件到存储卷
	file, err := os.Open(sourceFilePath)
	if err != nil {
		// 打开文件失败，删除卷并更新状态
		libvirtd.Conn.StorageVolDelete(vol, 0)
		image.Status = StatusError
		image.Save()
		return nil, fmt.Errorf("打开源文件失败: %v", err)
	}
	defer file.Close()

	// 上传文件到存储卷
	err = libvirtd.Conn.StorageVolUpload(vol, file, 0, fileSize, 0)
	if err != nil {
		// 上传失败，删除卷并更新状态
		libvirtd.Conn.StorageVolDelete(vol, 0)
		image.Status = StatusError
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
	image.Status = StatusActive
	if _, err := image.Save(); err != nil {
		return nil, fmt.Errorf("更新镜像状态失败: %v", err)
	}

	return image, nil
}

// TODO: 实现从URL创键镜像
// CreateImageFromURL 从URL创建镜像
func CreateImageFromURL(name, sourceURL, poolName, osType, arch, imageType, description, tags string, minDisk, minRam uint64) (*imageMod.Image, error) {
	// 这里可以实现从URL下载文件到临时目录，然后调用CreateImageFromLocalFile，最后清理临时文件
	// 由于涉及到网络下载，可能需要支持断点续传、进度显示等功能
	// 为了简化示例，这里省略具体实现
	return nil, fmt.Errorf("从URL创建镜像功能尚未实现")
}

// DeleteImage 删除镜像及关联的存储卷
func DeleteImage(idOrUUID string) error {
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
	image.Status = StatusDeleting
	if _, err := image.Save(); err != nil {
		return fmt.Errorf("更新镜像状态失败: %v", err)
	}

	// 3. 获取存储池
	pool, err := storagePool.GetStoragePool(image.PoolName)
	if err != nil {
		image.Status = StatusError
		image.Save()
		return fmt.Errorf("获取存储池失败: %v", err)
	}

	// 4. 删除存储卷
	err = volume.DeleteVolume(pool, image.VolumeName, 0)
	if err != nil {
		image.Status = StatusError
		image.Save()
		return fmt.Errorf("删除存储卷失败: %v", err)
	}

	// 5. 删除数据库记录
	if _, err := image.Delete(); err != nil {
		return fmt.Errorf("删除镜像记录失败: %v", err)
	}

	return nil
}
