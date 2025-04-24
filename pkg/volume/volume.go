package volume

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// ListVolumesSummary 列出存储池中的卷 简要信息
func ListVolumesSummary(Pool libvirt.StoragePool, Maxnames int32) (rNames []string, err error) {
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

// ListVolumesDetail 列出存储池中所有的卷 详细信息
func ListVolumesDetails(Pool libvirt.StoragePool, NeedResults int32, Flags uint32) (vols []libvirt.StorageVol, rRet uint32, err error) {
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

// CreateVolume 创建一个新的存储卷
func CreateVolume(Pool libvirt.StoragePool, Params *xmlDefine.VolumeTemplateParams, Flags libvirt.StorageVolCreateFlags) (vol libvirt.StorageVol, err error) {
	// 为所有未设置的字段应用默认值
	xmlDefine.SetDefaults(Params)

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.VolumeTemplate, Params)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("渲染XML失败: %w", err)
	}

	// 定义网络
	vol, err = libvirtd.Connection.StorageVolCreateXML(Pool, xmlStr, Flags)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("定义卷失败: %w", err)
	}

	return vol, nil
}
