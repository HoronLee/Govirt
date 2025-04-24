package volume

import (
	"fmt"
	"govirt/pkg/libvirtd"

	"github.com/digitalocean/go-libvirt"
)

// ListVolumes 列出存储池中的卷
func ListVolumes(Pool libvirt.StoragePool, Maxnames int32) (rNames []string, err error) {
	// 刷新存储池以确保获取最新信息
	err = libvirtd.Connection.StoragePoolRefresh(Pool, 0)
	if err != nil {
		return nil, fmt.Errorf("刷新存储池失败: %v", err)
	}

	// 获取存储池中的卷列表
	volumes, err := libvirtd.Connection.StoragePoolListVolumes(Pool, Maxnames)
	if err != nil {
		return nil, fmt.Errorf("列出存储池中的卷失败: %v", err)
	}
	return volumes, nil
}

// ListVolumesDetail 列出存储池中所有的卷附带详细信息
func ListVolumesDetail(Pool libvirt.StoragePool, NeedResults int32, Flags uint32) (vols []libvirt.StorageVol, rRet uint32, err error) {
	// 刷新存储池以确保获取最新信息
	err = libvirtd.Connection.StoragePoolRefresh(Pool, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("刷新存储池失败: %v", err)
	}

	// 获取存储池中的卷列表
	volumes, rRet, err := libvirtd.Connection.StoragePoolListAllVolumes(Pool, NeedResults, Flags)
	if err != nil {
		return nil, 0, fmt.Errorf("列出存储池中的卷失败: %v", err)
	}
	return volumes, rRet, nil
}

// GetVolumeInfo 获取特定卷的详细信息
func GetVolumeInfo(l *libvirt.Libvirt, PoolName, volumeName string) (rType int8, rCapacity uint64, rAllocation uint64, err error) {
	// 获取存储池
	Pool, err := libvirtd.Connection.StoragePoolLookupByName(PoolName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("查找存储池 %s 失败: %v", PoolName, err)
	}

	// 获取存储卷
	vol, err := libvirtd.Connection.StorageVolLookupByName(Pool, volumeName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("查找卷 %s 失败: %v", volumeName, err)
	}

	// 获取卷信息
	rType, rCapacity, rAllocation, err = libvirtd.Connection.StorageVolGetInfo(vol)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取卷 %s 信息失败: %v", volumeName, err)
	}

	return
}
