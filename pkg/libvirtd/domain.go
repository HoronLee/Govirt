// filepath: /home/horonlee/projects/govirt/pkg/libvirtd/domain.go
package libvirtd

import (
	"fmt"
	"govirt/pkg/config"
	"govirt/pkg/helpers"

	imageMod "govirt/app/models/image"
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
		Name: "Rocky9.2",
		// BootDev:      "cdrom",
		VCPU:       4,
		CurrentMem: 2097152,
		MaxMem:     2097152,
		// OsDiskSource: "/data/images/Rocky9.2Convert.qcow2",
		// CDRomSource:  "/data/images/Rocky-9.2-x86_64-minimal_174e5fe9-f573-96e3-aeeb-d40f4ab89bdf.iso",
		OsImageID: "Rocky9.2-Convert",
	}

	// 调用正式的创建方法
	_, err := vc.CreateDomainFromImage(params)
	if err != nil {
		logger.ErrorString("libvirt", "创建测试域失败", err.Error())
		return
	}

}

// CreateDomainFromImage 根据提供的参数从现有镜像克隆创建虚拟机
func (vc *VirtConn) CreateDomainFromImage(dparams *xmlDefine.DomainTemplateParams) (libvirt.Domain, error) {
	if dparams.OsImageID == "" {
		return libvirt.Domain{}, fmt.Errorf("必须提供 OsImageID 以便从镜像创建")
	}

	// --- 获取并检查基础镜像信息 ---
	imageID := dparams.OsImageID
	image, err := imageMod.GetByID(imageID)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("获取基础镜像失败: %w", err)
	}
	if image.Status != "active" {
		return libvirt.Domain{}, fmt.Errorf("基础镜像 '%s' 状态为 '%s'，不是 active 状态", image.Name, image.Status)
	}
	imagePool, err := vc.GetStoragePool(image.PoolName)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("获取基础镜像所在存储池失败: %w", err)
	}
	baseVolume, err := vc.GetVolume(imagePool, image.VolumeName)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("获取基础镜像存储卷失败: %w", err)
	}

	// --- 克隆卷作为新虚拟机的系统盘 ---
	osDiskPoolName := config.Get("pool.volume.name") // 系统盘目标存储池
	osDiskPool, err := vc.GetStoragePool(osDiskPoolName)
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("获取系统盘目标存储池 '%s' 失败: %w", osDiskPoolName, err)
	}
	// 定义新系统盘卷的参数
	osVolParams := &xmlDefine.VolumeTemplateParams{
		Name:     dparams.Name + ".qcow2",
		Capacity: dparams.OsCapacity, // 注意：克隆时容量通常由基础卷决定，除非CloneVolume支持调整
		Type:     "qcow2",
	}
	// 执行克隆
	osVolume, err := vc.CloneVolume(osDiskPool, osVolParams, baseVolume, 0) // flags=0
	if err != nil {
		return libvirt.Domain{}, fmt.Errorf("克隆基础镜像卷失败: %w", err)
	}
	// 设置系统盘源为新克隆的卷
	dparams.OsDiskSource = osVolume.Key // 使用克隆后卷的Key或Path

	xmlDefine.SetDefaults(dparams) // 应用默认值

	// 生成MAC地址
	if err := vc.ensureMacAddresses(dparams); err != nil {
		// 如果生成MAC失败，可能需要清理已克隆的卷
		_ = vc.DeleteVolume(osDiskPool, osVolume.Name, 0) // 尝试清理
		return libvirt.Domain{}, err
	}

	// 渲染XML模板
	xmlStr, err := xmlDefine.RenderTemplate(xmlDefine.DomainTemplate, dparams)
	if err != nil {
		// 清理已克隆的卷
		_ = vc.DeleteVolume(osDiskPool, osVolume.Name, 0) // 尝试清理
		return libvirt.Domain{}, fmt.Errorf("渲染域XML失败: %w", err)
	}

	// --- 定义域 ---
	domain, err := vc.Libvirt.DomainDefineXML(xmlStr)
	if err != nil {
		// 清理已克隆的卷
		_ = vc.DeleteVolume(osDiskPool, osVolume.Name, 0) // 尝试清理
		return libvirt.Domain{}, fmt.Errorf("定义域失败: %w", err)
	}

	logger.InfoString("libvirt", "成功从镜像定义域", dparams.Name)
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
