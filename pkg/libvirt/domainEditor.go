package libvirt

import (
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// GetDomainXMLDesc 获取指定域的XML描述
func GetDomainXMLDesc(domain libvirt.Domain) (string, error) {
	xmlDesc, err := Libvirt.DomainGetXMLDesc(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域XML描述失败", err.Error())
		return "", err
	}
	return xmlDesc, nil
}

// DefineDomain 定义并启动指定的域
func DefineDomain(xmlDesc string) (libvirt.Domain, error) {
	domain, err := Libvirt.DomainDefineXML(xmlDesc)
	if err != nil {
		logger.ErrorString("libvirt", "定义并启动域失败", err.Error())
		return domain, err
	}
	return domain, nil
}

// UpdateDomain 更新已存在的域定义
func UpdateDomain(domain libvirt.Domain, xmlDesc string) (libvirt.Domain, error) {
	// libvirt.DomainDefineXMLFlags 用于更新域定义，第二个参数是flags，0表示默认行为
	newDomain, err := Libvirt.DomainDefineXMLFlags(xmlDesc, 0)
	if err != nil {
		logger.ErrorString("libvirt", "更新域定义失败", err.Error())
		return libvirt.Domain{}, err
	}
	return newDomain, nil
}
