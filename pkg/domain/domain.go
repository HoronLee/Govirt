package domain

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// CreateATestDomain 创建一个测试虚拟机
// 这是一个演示如何使用CreateDomain函数的示例
func CreateATestDomain() {
	// 创建模板参数结构体实例
	// 只设置必要的参数，其他参数使用默认值
	params := &xmlDefine.DomainTemplateParams{
		Name: "test",
		// BootDev:      "cdrom",
		OsDiskSource: "/data/images/rocky9.qcow2",
		// CDRomSource:  "/data/images/Rocky-9.2-x86_64-minimal.iso",
	}

	// 调用正式的创建方法
	domain, err := CreateDomain(params)
	if err != nil {
		logger.ErrorString("libvirt", "创建测试域失败", err.Error())
		return
	}

	fmt.Printf("测试域创建成功: %v\n", domain.Name)

	// 输出完整参数，展示默认值被正确应用
	fmt.Printf("名称: %s\n", params.Name)
	fmt.Printf("UUID: %s\n", params.UUID)
	fmt.Printf("内存: %d KiB\n", params.MaxMem)
	fmt.Printf("当前内存: %d KiB\n", params.CurrentMem)
	fmt.Printf("VCPU: %d\n", params.VCPU)
	fmt.Printf("MAC地址: %s\n", params.NatMac)
}

// CreateDomain 根据提供的参数创建虚拟机
// 此函数接收一个已填充好的DomainTemplateParams结构体实例，用于创建新的虚拟机
func CreateDomain(params *xmlDefine.DomainTemplateParams) (libvirt.Domain, error) {
	// 为所有未设置的字段应用默认值
	xmlDefine.SetDefaults(params)

	// 如果未提供MAC地址，则自动生成一个
	if params.NatMac == "" {
		macAddr, err := helpers.GenerateRandomMAC()
		if err != nil {
			return libvirt.Domain{}, fmt.Errorf("生成随机MAC地址失败: %w", err)
		}
		params.NatMac = macAddr
	}

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.DomainTemplate, params)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("渲染域XML失败: %w", err)
	}

	// 定义域
	domain, err := libvirtd.Connection.DomainDefineXML(xmlStr)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("定义域失败: %w", err)
	}

	return domain, nil
}

// GetDomainXMLDesc 获取指定域的XML描述
func GetDomainXMLDesc(domain libvirt.Domain) (string, error) {
	xmlDesc, err := libvirtd.Connection.DomainGetXMLDesc(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域XML描述失败", err.Error())
		return "", err
	}
	return xmlDesc, nil
}

// DefineDomain 定义域
func DefineDomain(xmlDesc string) (libvirt.Domain, error) {
	domain, err := libvirtd.Connection.DomainDefineXML(xmlDesc)
	if err != nil {
		return domain, err
	}
	return domain, nil
}

// UpdateDomain 更新已存在的域定义
func UpdateDomain(domain libvirt.Domain, xmlDesc string) (libvirt.Domain, error) {
	// libvirt.DomainDefineXMLFlags 用于更新域定义，第二个参数是flags，0表示默认行为
	newDomain, err := libvirtd.Connection.DomainDefineXMLFlags(xmlDesc, 0)
	if err != nil {
		logger.ErrorString("libvirt", "更新域定义失败", err.Error())
		return libvirt.Domain{}, err
	}
	return newDomain, nil
}
