package volume

import (
	"fmt"
	"govirt/pkg/libvirtd"

	"github.com/digitalocean/go-libvirt"
)

// ListVolumes 列出存储池中的卷
func ListVolumes(pool libvirt.StoragePool, Maxnames int32) (rNames []string, err error) {
	// 刷新存储池以确保获取最新信息
	err = libvirtd.Connection.StoragePoolRefresh(pool, 0)
	if err != nil {
		return nil, fmt.Errorf("刷新存储池失败: %v", err)
	}

	// 获取存储池中的卷列表
	volumes, err := libvirtd.Connection.StoragePoolListVolumes(pool, Maxnames)
	if err != nil {
		return nil, fmt.Errorf("列出存储池中的卷失败: %v", err)
	}
	return volumes, nil
}

// GetVolumeInfo 获取特定卷的详细信息
func GetVolumeInfo(l *libvirt.Libvirt, poolName, volumeName string) (rType int8, rCapacity uint64, rAllocation uint64, err error) {
	// 获取存储池
	pool, err := libvirtd.Connection.StoragePoolLookupByName(poolName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("查找存储池 %s 失败: %v", poolName, err)
	}

	// 获取存储卷
	vol, err := libvirtd.Connection.StorageVolLookupByName(pool, volumeName)
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
