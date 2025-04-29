// filepath: /home/horonlee/projects/govirt/pkg/libvirtd/domain.go
package libvirtd

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// CreateATestDomain 创建一个测试虚拟机
// 这是一个演示如何使用CreateDomain函数的示例
// TODO: 实现创建domian的逻辑之前需要完成镜像管理
func (vc *VirtConn) CreateATestDomain() {
	// 创建模板参数结构体实例
	// 只设置必要的参数，其他参数使用默认值
	params := &xmlDefine.DomainTemplateParams{
		Name: "test",
		// BootDev:      "cdrom",
		OsDiskSource: "/data/images/rocky9.qcow2",
		// CDRomSource:  "/data/images/Rocky-9.2-x86_64-minimal.iso",
	}

	// 调用正式的创建方法
	domain, err := vc.CreateDomain(params)
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
func (vc *VirtConn) CreateDomain(params *xmlDefine.DomainTemplateParams) (libvirt.Domain, error) {
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
	domain, err := vc.Libvirt.DomainDefineXML(xmlStr)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("定义域失败: %w", err)
	}

	return domain, nil
}

// GetDomainXMLDesc 获取指定域的XML描述
func (vc *VirtConn) GetDomainXMLDesc(domain libvirt.Domain) (string, error) {
	xmlDesc, err := vc.Libvirt.DomainGetXMLDesc(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域XML描述失败", err.Error())
		return "", err
	}
	return xmlDesc, nil
}

// DefineDomain 定义域
func (vc *VirtConn) DefineDomain(xmlDesc string) (libvirt.Domain, error) {
	domain, err := vc.Libvirt.DomainDefineXML(xmlDesc)
	if err != nil {
		return domain, err
	}
	return domain, nil
}

// UpdateDomain 更新已存在的域定义
func (vc *VirtConn) UpdateDomain(domain libvirt.Domain, xmlDesc string) (libvirt.Domain, error) {
	// libvirt.DomainDefineXMLFlags 用于更新域定义，第二个参数是flags，0表示默认行为
	newDomain, err := vc.Libvirt.DomainDefineXMLFlags(xmlDesc, 0)
	if err != nil {
		logger.ErrorString("libvirt", "更新域定义失败", err.Error())
		return libvirt.Domain{}, err
	}
	return newDomain, nil
}

// GetDomain 根据 UUID、UUID字符串或名称获取域
func (vc *VirtConn) GetDomain(identifier any) (libvirt.Domain, error) {
	domains, err := vc.ListAllDomains()
	if err != nil {
		logger.ErrorString("libvirt", "获取所有域失败", err.Error())
		return libvirt.Domain{}, err
	}

	for _, domain := range domains {
		switch id := identifier.(type) {
		case libvirt.UUID:
			if domain.UUID == id {
				return domain, nil
			}
		case string:
			// 检查是否为UUID字符串格式
			if helpers.IsUUIDString(id) {
				uuid, err := helpers.UUIDStringToBytes(id)
				if err != nil {
					continue
				}
				if domain.UUID == uuid {
					return domain, nil
				}
			} else if domain.Name == id {
				// 如果不是UUID格式，则按名称匹配
				return domain, nil
			}
		default:
			return libvirt.Domain{}, fmt.Errorf("无效的标识符类型: %T", identifier)
		}
	}

	return libvirt.Domain{}, nil
}

// UpdateDomainStateByUUID 根据 UUID 更新域的操作状态
func (vc *VirtConn) UpdateDomainStateByUUID(uuid libvirt.UUID, op DomainOperation, flags uint32) (libvirt.DomainState, error) {
	domain, err := vc.GetDomain(uuid)
	if err != nil {
		return libvirt.DomainNostate, err
	}
	currentState, err := vc.GetDomainState(domain)
	if err != nil {
		return currentState, err
	}
	switch op {
	case DomainOpStart:
		if currentState != libvirt.DomainRunning {
			err = vc.StartDomain(domain)
		}
	case DomainOpShutdown:
		if currentState == libvirt.DomainRunning {
			err = vc.ShutdownDomain(domain)
		}
	case DomainOpForceStop:
		if currentState == libvirt.DomainRunning || currentState == libvirt.DomainPaused {
			err = vc.ForceStopDomain(domain)
		}
	case DomainOpReboot:
		if currentState == libvirt.DomainRunning {
			err = vc.RebootDomain(domain)
		}
	case DomainOpForceReboot:
		if currentState == libvirt.DomainRunning {
			err = vc.ForceRebootDomain(domain)
		}
	case DomainOpSuspend:
		if currentState == libvirt.DomainRunning {
			err = vc.SuspendDomain(domain)
		}
	case DomainOpResume:
		if currentState == libvirt.DomainPaused {
			err = vc.ResumeDomain(domain)
		}
	case DomainOpSave:
		if currentState == libvirt.DomainRunning {
			err = vc.SaveDomain(domain)
		}
	case DomainOpRestore:
		err = fmt.Errorf("恢复操作暂未实现")
	case DomainOpDelete:
		// 注意：这里调用 ForceDeleteDomain，因为它包含了停止逻辑
		// 原来的 flags 参数在这里不再需要，因为 ForceDeleteDomain 内部处理了
		err = vc.ForceDeleteDomain(domain)
	case DomainOpClone, DomainOpMigrate, DomainOpSnapshot:
		err = fmt.Errorf("该操作暂未实现")
	default:
		return currentState, fmt.Errorf("未知操作")
	}
	if err != nil {
		logger.ErrorString("libvirt", "执行域操作失败", err.Error())
		return currentState, err
	}
	return vc.GetDomainState(domain)
}

// ListAllDomains 列出所有域的信息
func (vc *VirtConn) ListAllDomains() ([]libvirt.Domain, error) {
	domains, _, err := vc.Libvirt.ConnectListAllDomains(1, 0)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有域失败", err.Error())
		return nil, err
	}
	return domains, nil
}

// GetDomainState 获取指定域的状态
func (vc *VirtConn) GetDomainState(domain libvirt.Domain) (libvirt.DomainState, error) {
	state, _, err := vc.Libvirt.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// GetDomainStateByUUID 根据 UUID 获取域的状态
func (vc *VirtConn) GetDomainStateByUUID(uuid libvirt.UUID) (libvirt.DomainState, error) {
	domain, err := vc.GetDomain(uuid)
	if err != nil {
		return libvirt.DomainState(0), err
	}
	state, _, err := vc.Libvirt.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// StartDomain 开机
func (vc *VirtConn) StartDomain(domain libvirt.Domain) error {
	err := vc.Libvirt.DomainCreate(domain)
	if err != nil {
		logger.ErrorString("libvirt", "启动域失败", err.Error())
		return err
	}
	return nil
}

// ShutdownDomain 正常关机
func (vc *VirtConn) ShutdownDomain(domain libvirt.Domain) error {
	err := vc.Libvirt.DomainShutdownFlags(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "关闭域失败", err.Error())
		return err
	}
	return nil
}

// ForceStopDomain 强制关机
func (vc *VirtConn) ForceStopDomain(domain libvirt.Domain) error {
	err := vc.Libvirt.DomainDestroyFlags(domain, libvirt.DomainDestroyDefault)
	if err != nil {
		logger.ErrorString("libvirt", "强制停止域失败", err.Error())
		return err
	}
	return nil
}

// SuspendDomain 暂停
func (vc *VirtConn) SuspendDomain(domain libvirt.Domain) error {
	err := vc.Libvirt.DomainSuspend(domain)
	if err != nil {
		logger.ErrorString("libvirt", "暂停域失败", err.Error())
		return err
	}
	return nil
}

// ResumeDomain 恢复
func (vc *VirtConn) ResumeDomain(domain libvirt.Domain) error {
	err := vc.Libvirt.DomainResume(domain)
	if err != nil {
		logger.ErrorString("libvirt", "恢复域失败", err.Error())
		return err
	}
	return nil
}

// RebootDomain 重启
func (vc *VirtConn) RebootDomain(domain libvirt.Domain) error {
	if err := vc.ShutdownDomain(domain); err != nil {
		return err
	}
	return vc.StartDomain(domain)
}

// ForceRebootDomain 强制重启
func (vc *VirtConn) ForceRebootDomain(domain libvirt.Domain) error {
	if err := vc.ForceStopDomain(domain); err != nil {
		return err
	}
	return vc.StartDomain(domain)
}

// SaveDomain 保存状态
func (vc *VirtConn) SaveDomain(domain libvirt.Domain) error {
	// 自动生成保存路径：/var/lib/libvirt/save/<domain-name>.save
	savePath := fmt.Sprintf("/var/lib/libvirt/save/%s.save", domain.Name)
	err := vc.Libvirt.DomainSave(domain, savePath)
	if err != nil {
		logger.ErrorString("libvirt", "保存域状态失败", err.Error())
		return err
	}
	return nil
}

// ForceDeleteDomain 强制删除指定的域（会先强制停止）
// 会删除快照元数据和NVRAM配置
func (vc *VirtConn) ForceDeleteDomain(domain libvirt.Domain) error {
	currentState, err := vc.GetDomainState(domain)
	if err != nil {
		// 如果获取状态失败，仍然尝试删除定义，可能域已经损坏
		logger.WarnString("libvirt", "获取域状态失败，仍尝试强制删除", err.Error())
	} else {
		// 如果域正在运行或暂停，则强制停止
		if currentState == libvirt.DomainRunning || currentState == libvirt.DomainPaused {
			logger.InfoString("libvirt", "域正在运行或暂停，执行强制停止", domain.Name)
			if err := vc.ForceStopDomain(domain); err != nil {
				logger.ErrorString("libvirt", "强制停止域失败，无法继续删除", err.Error())
				return fmt.Errorf("强制停止域失败: %w", err)
			}
			// 强制停止后需要一点时间让状态更新，或者依赖后续UndefineFlags能处理
		}
	}

	// 使用包含快照和NVRAM的删除标志
	flags := libvirt.DomainUndefineSnapshotsMetadata | libvirt.DomainUndefineNvram | libvirt.DomainUndefineManagedSave | libvirt.DomainUndefineCheckpointsMetadata
	err = vc.Libvirt.DomainUndefineFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "强制删除域定义失败", err.Error())
		return err
	}
	logger.InfoString("libvirt", "成功强制删除域", domain.Name)
	return nil
}

// DeleteStoppedDomain 删除已停止的域
// 会删除快照元数据和NVRAM配置
func (vc *VirtConn) DeleteStoppedDomain(domain libvirt.Domain) error {
	currentState, err := vc.GetDomainState(domain)
	if err != nil {
		// 如果获取状态失败，也允许尝试删除，可能域已损坏
		logger.WarnString("libvirt", "获取域状态失败，仍尝试删除已停止的域", err.Error())
	} else {
		// 检查域是否处于非运行状态
		switch currentState {
		case libvirt.DomainShutoff, libvirt.DomainShutdown, libvirt.DomainCrashed:
			// 状态符合要求，继续执行删除
		case libvirt.DomainRunning, libvirt.DomainPaused, libvirt.DomainBlocked, libvirt.DomainPmsuspended:
			// 状态不符合要求，返回错误
			err := fmt.Errorf("域 '%s' 处于活动状态 (%v)，无法删除，请先停止", domain.Name, currentState)
			logger.ErrorString("libvirt", "删除已停止的域失败", err.Error())
			return err
		default: // DomainNostate 或其他未知状态
			logger.WarnString("libvirt", fmt.Sprintf("域 '%s' 状态未知 (%v)，尝试删除", domain.Name, currentState), "")
		}
	}

	// 使用包含快照和NVRAM的删除标志 (与原默认行为一致)
	flags := libvirt.DomainUndefineSnapshotsMetadata | libvirt.DomainUndefineNvram
	err = vc.Libvirt.DomainUndefineFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "删除已停止的域定义失败", err.Error())
		return err
	}
	logger.InfoString("libvirt", "成功删除已停止的域", domain.Name)
	return nil
}

// SetDomainAutostart 设置虚拟机自动启动
func (vc *VirtConn) SetDomainAutostart(domain libvirt.Domain, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := vc.Libvirt.DomainSetAutostart(domain, autostartFlag); err != nil {
		return fmt.Errorf("设置虚拟机自动启动失败: %w", err)
	}
	return nil
}

// GetDomainAutostart 获取虚拟机自动启动状态
func (vc *VirtConn) GetDomainAutostart(domain libvirt.Domain) (bool, error) {
	autostart, err := vc.Libvirt.DomainGetAutostart(domain)
	if err != nil {
		return false, fmt.Errorf("获取虚拟机自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
