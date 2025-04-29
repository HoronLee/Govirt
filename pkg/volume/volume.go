package volume

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/xmlDefine"
	"io"

	"github.com/digitalocean/go-libvirt"
)

// ListVolumesSummary 列出存储池中的卷 简要信息
func ListVolumesSummary(Pool libvirt.StoragePool) (rNames []string, err error) {
	// 刷新存储池以确保获取最新信息
	err = libvirtd.Conn.StoragePoolRefresh(Pool, 0)
	if err != nil {
		return nil, fmt.Errorf("刷新存储池失败: %v", err)
	}
	// 获取存储池中的卷数量
	resultNum, err := GetVolumeNum(Pool)
	if err != nil {
		return nil, fmt.Errorf("获取存储池卷数量失败: %v", err)
	}
	// 获取存储池中的卷列表
	volumes, err := libvirtd.Conn.StoragePoolListVolumes(Pool, resultNum)
	if err != nil {
		return nil, fmt.Errorf("列出存储池中的卷失败: %v", err)
	}
	return volumes, nil
}

// ListVolumesDetail 列出存储池中所有的卷 详细信息
func ListVolumesDetails(Pool libvirt.StoragePool, Flags uint32) (vols []libvirt.StorageVol, rRet uint32, err error) {
	// 刷新存储池以确保获取最新信息
	err = libvirtd.Conn.StoragePoolRefresh(Pool, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("刷新存储池失败: %v", err)
	}

	// 获取存储池中的卷列表
	volumes, rRet, err := libvirtd.Conn.StoragePoolListAllVolumes(Pool, 1, Flags)
	if err != nil {
		return nil, 0, fmt.Errorf("列出存储池中的卷失败: %v", err)
	}
	return volumes, rRet, nil
}

// CreateVolume 创建一个新的存储卷
func CreateVolume(Pool libvirt.StoragePool, Params *xmlDefine.VolumeTemplateParams, Flags libvirt.StorageVolCreateFlags) (vol libvirt.StorageVol, err error) {
	xmlDefine.SetDefaults(Params)

	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.VolumeTemplate, Params)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	vol, err = libvirtd.Conn.StorageVolCreateXML(Pool, xmlStr, Flags)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("定义卷失败: %w", err)
	}

	return vol, nil
}

// DeleteVolume 删除存储卷
func DeleteVolume(Pool libvirt.StoragePool, VolumeName string, Flags libvirt.StorageVolDeleteFlags) (err error) {
	// 获取存储卷
	vol, err := GetVolume(Pool, VolumeName)
	if err != nil {
		return fmt.Errorf("查找卷 %s 失败: %v", VolumeName, err)
	}

	// 删除存储卷
	err = libvirtd.Conn.StorageVolDelete(vol, Flags)
	if err != nil {
		return fmt.Errorf("删除卷 %s 失败: %v", VolumeName, err)
	}

	return nil
}

// CloneVolume 从现有存储卷克隆创建新的存储卷
func CloneVolume(Pool libvirt.StoragePool, NewParams *xmlDefine.VolumeTemplateParams, SourceVol libvirt.StorageVol, Flags libvirt.StorageVolCreateFlags) (vol libvirt.StorageVol, err error) {
	xmlDefine.SetDefaults(NewParams)

	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.VolumeTemplate, NewParams)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	vol, err = libvirtd.Conn.StorageVolCreateXMLFrom(Pool, xmlStr, SourceVol, Flags)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("克隆卷失败: %w", err)
	}

	return vol, nil
}

// CloneVolume 克隆存储卷 弃用的方法
func CloneVolumeByPipe(SourcePool libvirt.StoragePool, SourceVolName string,
	DestPool libvirt.StoragePool, NewParams *xmlDefine.VolumeTemplateParams,
	Flags libvirt.StorageVolCreateFlags) (vol libvirt.StorageVol, err error) {
	// 1. 获取源卷
	sourceVol, err := GetVolume(SourcePool, SourceVolName)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("获取源卷失败: %w", err)
	}

	// 获取源卷信息
	_, rCapacity, _, err := libvirtd.Conn.StorageVolGetInfo(sourceVol)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("获取源卷信息失败: %w", err)
	}

	// 确保目标卷参数中设置了正确的容量
	if NewParams.Capacity == 0 {
		NewParams.Capacity = rCapacity
	}

	// 在目标池中创建新卷
	newVol, err := CreateVolume(DestPool, NewParams, Flags)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("在目标池中创建卷失败: %w", err)
	}

	// 复制源卷数据到新卷
	err = copyVolData(sourceVol, newVol)
	if err != nil {
		// 如果复制失败，尝试删除目标卷
		_ = libvirtd.Conn.StorageVolDelete(newVol, 0)
		return libvirt.StorageVol{}, fmt.Errorf("复制卷数据失败: %w", err)
	}

	return newVol, nil
}

// copyVolData 在两个卷之间复制数据 弃用的方法
func copyVolData(sourceVol, destVol libvirt.StorageVol) error {
	// 获取源卷的大小信息
	_, capacity, _, err := libvirtd.Conn.StorageVolGetInfo(sourceVol)
	if err != nil {
		return fmt.Errorf("获取源卷信息失败: %w", err)
	}

	// 创建一个管道用于数据传输
	pipeReader, pipeWriter := io.Pipe()

	// 创建错误通道接收异步操作的错误
	errChan := make(chan error, 2)

	// 启动下载协程
	go func() {
		defer pipeWriter.Close()
		// StorageVolDownload 直接接收 Writer，将卷数据写入 pipeWriter
		err := libvirtd.Conn.StorageVolDownload(sourceVol, pipeWriter, 0, capacity, 0)
		errChan <- err
	}()

	// 启动上传协程
	go func() {
		defer pipeReader.Close()
		// StorageVolUpload 直接接收 Reader，从 pipeReader 读取数据
		err := libvirtd.Conn.StorageVolUpload(destVol, pipeReader, 0, capacity, 0)
		errChan <- err
	}()

	// 等待两个操作完成并检查错误
	var downloadErr, uploadErr error
	downloadErr = <-errChan
	uploadErr = <-errChan

	if downloadErr != nil {
		return fmt.Errorf("从源卷下载失败: %w", downloadErr)
	}
	if uploadErr != nil {
		return fmt.Errorf("上传到目标卷失败: %w", uploadErr)
	}

	return nil
}
