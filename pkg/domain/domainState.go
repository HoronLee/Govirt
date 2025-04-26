package domain

import (
	"fmt"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// UpdateDomainStateByUUID 根据 UUID 更新域的操作状态
func UpdateDomainStateByUUID(uuid libvirt.UUID, op DomainOperation, flags uint32) (libvirt.DomainState, error) {
	domain, err := GetDomain(uuid)
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
			err = ShutdownDomain(domain, libvirt.DomainShutdownFlagValues(flags))
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
		// 注意：这里调用 ForceDeleteDomain，因为它包含了停止逻辑
		// 原来的 flags 参数在这里不再需要，因为 ForceDeleteDomain 内部处理了
		err = ForceDeleteDomain(domain)
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
	domains, _, err := libvirtd.Connection.ConnectListAllDomains(1, 1|2)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有域失败", err.Error())
		return nil, err
	}
	return domains, nil
}

// GetDomainState 获取指定域的状态
func GetDomainState(domain libvirt.Domain) (libvirt.DomainState, error) {
	state, _, err := libvirtd.Connection.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// GetDomainStateByUUID 根据 UUID 获取域的状态
func GetDomainStateByUUID(uuid libvirt.UUID) (libvirt.DomainState, error) {
	domain, err := GetDomain(uuid)
	if err != nil {
		return libvirt.DomainState(0), err
	}
	state, _, err := libvirtd.Connection.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// StartDomain 开机
func StartDomain(domain libvirt.Domain) error {
	err := libvirtd.Connection.DomainCreate(domain)
	if err != nil {
		logger.ErrorString("libvirt", "启动域失败", err.Error())
		return err
	}
	return nil
}

// 枚举常量	                 	 值  作用
// DomainShutdownDefault	    0	默认关机方式，由 hypervisor 自动选择合适的关机方式
// DomainShutdownAcpiPowerBtn	1	模拟按下 ACPI 电源按钮，向虚拟机发送 ACPI 关机事件
// DomainShutdownGuestAgent	    2	通过客户机代理(guest agent)关闭，前提是在虚拟机中安装了 guest agent
// DomainShutdownInitctl		4	通过 initctl 机制关闭客户机，主要用于老的 Linux 系统
// DomainShutdownSignal			8	发送信号（通常是 SIGTERM）给虚拟机的初始进程
// DomainShutdownParavirt		16	使用半虚拟化关机接口，适用于支持半虚拟化的客户机

// ShutdownDomain 正常关机
func ShutdownDomain(domain libvirt.Domain, flags libvirt.DomainShutdownFlagValues) error {
	err := libvirtd.Connection.DomainShutdownFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "关闭域失败", err.Error())
		return err
	}
	return nil
}

// ForceStopDomain 强制关机
func ForceStopDomain(domain libvirt.Domain) error {
	err := libvirtd.Connection.DomainDestroyFlags(domain, libvirt.DomainDestroyDefault)
	if err != nil {
		logger.ErrorString("libvirt", "强制停止域失败", err.Error())
		return err
	}
	return nil
}

// SuspendDomain 暂停
func SuspendDomain(domain libvirt.Domain) error {
	err := libvirtd.Connection.DomainSuspend(domain)
	if err != nil {
		logger.ErrorString("libvirt", "暂停域失败", err.Error())
		return err
	}
	return nil
}

// ResumeDomain 恢复
func ResumeDomain(domain libvirt.Domain) error {
	err := libvirtd.Connection.DomainResume(domain)
	if err != nil {
		logger.ErrorString("libvirt", "恢复域失败", err.Error())
		return err
	}
	return nil
}

// RebootDomain 重启
func RebootDomain(domain libvirt.Domain) error {
	if err := ShutdownDomain(domain, 0); err != nil {
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
	err := libvirtd.Connection.DomainSave(domain, savePath)
	if err != nil {
		logger.ErrorString("libvirt", "保存域状态失败", err.Error())
		return err
	}
	return nil
}

// ForceDeleteDomain 强制删除指定的域（会先强制停止）
// 会删除快照元数据和NVRAM配置
func ForceDeleteDomain(domain libvirt.Domain) error {
	currentState, err := GetDomainState(domain)
	if err != nil {
		// 如果获取状态失败，仍然尝试删除定义，可能域已经损坏
		logger.WarnString("libvirt", "获取域状态失败，仍尝试强制删除", err.Error())
	} else {
		// 如果域正在运行或暂停，则强制停止
		if currentState == libvirt.DomainRunning || currentState == libvirt.DomainPaused {
			logger.InfoString("libvirt", "域正在运行或暂停，执行强制停止", domain.Name)
			if err := ForceStopDomain(domain); err != nil {
				logger.ErrorString("libvirt", "强制停止域失败，无法继续删除", err.Error())
				return fmt.Errorf("强制停止域失败: %w", err)
			}
			// 强制停止后需要一点时间让状态更新，或者依赖后续UndefineFlags能处理
		}
	}

	// 使用包含快照和NVRAM的删除标志
	flags := libvirt.DomainUndefineSnapshotsMetadata | libvirt.DomainUndefineNvram | libvirt.DomainUndefineManagedSave | libvirt.DomainUndefineCheckpointsMetadata
	err = libvirtd.Connection.DomainUndefineFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "强制删除域定义失败", err.Error())
		return err
	}
	logger.InfoString("libvirt", "成功强制删除域", domain.Name)
	return nil
}

// DeleteStoppedDomain 删除已停止的域
// 会删除快照元数据和NVRAM配置
func DeleteStoppedDomain(domain libvirt.Domain) error {
	currentState, err := GetDomainState(domain)
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
	err = libvirtd.Connection.DomainUndefineFlags(domain, flags)
	if err != nil {
		logger.ErrorString("libvirt", "删除已停止的域定义失败", err.Error())
		return err
	}
	logger.InfoString("libvirt", "成功删除已停止的域", domain.Name)
	return nil
}

// SetDomainAutostart 设置虚拟机自动启动
func SetDomainAutostart(domain libvirt.Domain, autostart bool) error {
	var autostartFlag int32 = 0
	if autostart {
		autostartFlag = 1
	}

	if err := libvirtd.Connection.DomainSetAutostart(domain, autostartFlag); err != nil {
		return fmt.Errorf("设置虚拟机自动启动失败: %w", err)
	}
	return nil
}

// GetDomainAutostart 获取虚拟机自动启动状态
func GetDomainAutostart(domain libvirt.Domain) (bool, error) {
	autostart, err := libvirtd.Connection.DomainGetAutostart(domain)
	if err != nil {
		return false, fmt.Errorf("获取虚拟机自动启动状态失败: %w", err)
	}

	return autostart == 1, nil
}
