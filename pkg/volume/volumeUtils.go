package volume

import (
	"fmt"
	"govirt/pkg/libvirtd"

	"github.com/digitalocean/go-libvirt"
)

// GetVolume 获取卷
func GetVolume(Pool libvirt.StoragePool, VolumeName string) (vol libvirt.StorageVol, err error) {
	// 获取存储卷
	vol, err = libvirtd.Connection.StorageVolLookupByName(Pool, VolumeName)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("查找卷 %s 失败: %v", VolumeName, err)
	}

	return vol, nil
}

// GetVolumeInfo 获取特定卷的详细信息
func GetVolumeInfo(Pool libvirt.StoragePool, VolumeName string) (rType int8, rCapacity uint64, rAllocation uint64, err error) {
	// 获取存储卷
	vol, err := GetVolume(Pool, VolumeName)
	if err != nil {
		return 0, 0, 0, err
	}

	// 获取卷信息
	rType, rCapacity, rAllocation, err = libvirtd.Connection.StorageVolGetInfo(vol)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("获取卷 %s 信息失败: %v", VolumeName, err)
	}

	return
}

// GetVolumeNum 获取存储池中的卷数量
func GetVolumeNum(Pool libvirt.StoragePool) (rNum int32, err error) {
	// 获取存储池中的卷数量
	rNum, err = libvirtd.Connection.StoragePoolNumOfVolumes(Pool)
	if err != nil {
		return 0, fmt.Errorf("获取存储池 %s 中的卷数量失败: %v", Pool.Name, err)
	}

	return
}
