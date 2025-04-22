package libvirt

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
)

// CreateATestDomain 创建一个测试虚拟机
// 这是一个演示如何使用CreateDomain函数的示例
func CreateATestDomain() {
	// 创建模板参数结构体实例
	// 只设置必要的参数，其他参数使用默认值
	params := &DomainTemplateParams{
		Name:          "test",
		OsDiskSource:  "/data/images/rocky9.qcow2",
		VncPort:       "-1",
		IsVncAutoPort: "no",
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
func CreateDomain(params *DomainTemplateParams) (libvirt.Domain, error) {
	// 为所有未设置的字段应用默认值
	SetDefaults(params)

	// 如果未提供MAC地址，则自动生成一个
	if params.NatMac == "" {
		macAddr, err := helpers.GenerateRandomMAC()
		if err != nil {
			return libvirt.Domain{}, fmt.Errorf("生成随机MAC地址失败: %w", err)
		}
		params.NatMac = macAddr
	}

	// 如果未提供UUID，则自动生成一个
	if params.UUID == "" {
		params.UUID = uuid.New().String()
	}

	// 渲染XML模板
	xmlStr, err := RenderTemplate(domainTemplate, params)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("渲染域XML失败: %w", err)
	}

	// 定义域
	domain, err := connection.DomainDefineXML(xmlStr)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("定义域失败: %w", err)
	}

	return domain, nil
}

// GetDomainXMLDesc 获取指定域的XML描述
func GetDomainXMLDesc(domain libvirt.Domain) (string, error) {
	xmlDesc, err := connection.DomainGetXMLDesc(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域XML描述失败", err.Error())
		return "", err
	}
	return xmlDesc, nil
}

// DefineDomain 定义域
func DefineDomain(xmlDesc string) (libvirt.Domain, error) {
	domain, err := connection.DomainDefineXML(xmlDesc)
	if err != nil {
		return domain, err
	}
	return domain, nil
}

// UpdateDomain 更新已存在的域定义
func UpdateDomain(domain libvirt.Domain, xmlDesc string) (libvirt.Domain, error) {
	// libvirt.DomainDefineXMLFlags 用于更新域定义，第二个参数是flags，0表示默认行为
	newDomain, err := connection.DomainDefineXMLFlags(xmlDesc, 0)
	if err != nil {
		logger.ErrorString("libvirt", "更新域定义失败", err.Error())
		return libvirt.Domain{}, err
	}
	return newDomain, nil
}

// UpdateDomainStateByUUID 根据 UUID 更新域的操作状态
func UpdateDomainStateByUUID(uuid libvirt.UUID, op DomainOperation, flag libvirt.DomainUndefineFlagsValues) (libvirt.DomainState, error) {
	domain, err := GetDomainByUUID(uuid)
	if err != nil {
		return libvirt.DomainNostate, err
	}
	currentState, err := GetDomainState(domain)
	if err != nil {
		return currentState, err
	}
	switch op {
	case DomainOpStart:
		if currentState != libvirt.DomainRunning {
			err = StartDomain(domain)
		}
	case DomainOpShutdown:
		if currentState == libvirt.DomainRunning {
			err = ShutdownDomain(domain)
		}
	case DomainOpForceStop:
		if currentState == libvirt.DomainRunning || currentState == libvirt.DomainPaused {
			err = ForceStopDomain(domain)
		}
	case DomainOpReboot:
		if currentState == libvirt.DomainRunning {
			err = RebootDomain(domain)
		}
	case DomainOpForceReboot:
		if currentState == libvirt.DomainRunning {
			err = ForceRebootDomain(domain)
		}
	case DomainOpSuspend:
		if currentState == libvirt.DomainRunning {
			err = SuspendDomain(domain)
		}
	case DomainOpResume:
		if currentState == libvirt.DomainPaused {
			err = ResumeDomain(domain)
		}
	case DomainOpSave:
		if currentState == libvirt.DomainRunning {
			err = SaveDomain(domain)
		}
	case DomainOpRestore:
		err = fmt.Errorf("恢复操作暂未实现")
	case DomainOpDelete:
		if currentState != libvirt.DomainRunning {
			err = DeleteDomain(domain, flag) // 使用默认删除标志
		} else {
			err = fmt.Errorf("无法删除运行中的域")
		}
	case DomainOpClone, DomainOpMigrate, DomainOpSnapshot:
		err = fmt.Errorf("该操作暂未实现")
	default:
		return currentState, fmt.Errorf("未知操作")
	}
	if err != nil {
		logger.ErrorString("libvirt", "执行域操作失败", err.Error())
		return currentState, err
	}
	return GetDomainState(domain)
}

// ListAllDomains 列出所有域的信息
func ListAllDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := connection.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有域失败", err.Error())
		return nil, err
	}
	return domains, nil
}

// ListActiveDomains 列出所有活动域的信息
func ListActiveDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsActive
	domains, _, err := connection.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有活动域失败", err.Error())
	}
	return domains, nil
}

// ListInactiveDomains 列出所有非活动域的信息
func ListInactiveDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsInactive
	domains, _, err := connection.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有非活动域失败", err.Error())
	}
	return domains, nil
}

// GetDomainState 获取指定域的状态
func GetDomainState(domain libvirt.Domain) (libvirt.DomainState, error) {
	state, _, err := connection.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// GetDomainStateByUUID 根据 UUID 获取域的状态
func GetDomainStateByUUID(uuid libvirt.UUID) (libvirt.DomainState, error) {
	domain, err := GetDomainByUUID(uuid)
	if err != nil {
		return libvirt.DomainState(0), err
	}
	state, _, err := connection.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// StartDomain 开机
func StartDomain(domain libvirt.Domain) error {
	err := connection.DomainCreate(domain)
	if err != nil {
		logger.ErrorString("libvirt", "启动域失败", err.Error())
		return err
	}
	return nil
}

// ShutdownDomain 正常关机
func ShutdownDomain(domain libvirt.Domain) error {
	err := connection.DomainShutdownFlags(domain, libvirt.DomainShutdownDefault)
	if err != nil {
		logger.ErrorString("libvirt", "关闭域失败", err.Error())
		return err
	}
	return nil
}

// ForceStopDomain 强制关机
func ForceStopDomain(domain libvirt.Domain) error {
	err := connection.DomainDestroyFlags(domain, libvirt.DomainDestroyDefault)
	if err != nil {
		logger.ErrorString("libvirt", "强制停止域失败", err.Error())
		return err
	}
	return nil
}

// SuspendDomain 暂停
func SuspendDomain(domain libvirt.Domain) error {
	err := connection.DomainSuspend(domain)
	if err != nil {
		logger.ErrorString("libvirt", "暂停域失败", err.Error())
		return err
	}
	return nil
}

// ResumeDomain 恢复
func ResumeDomain(domain libvirt.Domain) error {
	err := connection.DomainResume(domain)
	if err != nil {
		logger.ErrorString("libvirt", "恢复域失败", err.Error())
		return err
	}
	return nil
}

// RebootDomain 重启
func RebootDomain(domain libvirt.Domain) error {
	if err := ShutdownDomain(domain); err != nil {
		return err
	}
	return StartDomain(domain)
}

// ForceRebootDomain 强制重启
func ForceRebootDomain(domain libvirt.Domain) error {
	if err := ForceStopDomain(domain); err != nil {
		return err
	}
	return StartDomain(domain)
}

// SaveDomain 保存状态
func SaveDomain(domain libvirt.Domain) error {
	// 自动生成保存路径：/var/lib/libvirt/save/<domain-name>.save
	savePath := fmt.Sprintf("/var/lib/libvirt/save/%s.save", domain.Name)
	err := connection.DomainSave(domain, savePath)
	if err != nil {
		logger.ErrorString("libvirt", "保存域状态失败", err.Error())
		return err
	}
	return nil
}

// DeleteDomain 删除指定的域
// flags 为0时使用默认删除标志(删除快照元数据和NVRAM配置)
func DeleteDomain(domain libvirt.Domain, flags libvirt.DomainUndefineFlagsValues) error {
	// 如果未指定flags，使用默认删除标志
	if flags == 0 {
		flags = libvirt.DomainUndefineSnapshotsMetadata | libvirt.DomainUndefineNvram
	}

	err := connection.DomainUndefineFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "删除域失败", err.Error())
		return err
	}
	return nil
}
