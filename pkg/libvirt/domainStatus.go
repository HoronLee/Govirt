package libvirt

import (
	"gohub/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// ListAllDomains 列出所有域的信息
func ListAllDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := Libvirt.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有域失败", err.Error())
		return nil, err
	}
	return domains, nil
}

// ListActiveDomains 列出所有活动域的信息
func ListActiveDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsActive
	domains, _, err := Libvirt.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有活动域失败", err.Error())
	}
	return domains, nil
}

// ListInactiveDomains 列出所有非活动域的信息
func ListInactiveDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsInactive
	domains, _, err := Libvirt.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有非活动域失败", err.Error())
	}
	return domains, nil
}

// GetDomainState 获取指定域的状态
func GetDomainState(domain libvirt.Domain) (libvirt.DomainState, error) {
	state, _, err := Libvirt.DomainGetState(domain, 0)
	if err != nil {
		logger.ErrorString("libvirt", "获取域状态失败", err.Error())
		return libvirt.DomainState(state), err
	}
	return libvirt.DomainState(state), nil
}

// StartDomain 启动指定的域
func StartDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainCreate(domain)
	if err != nil {
		logger.ErrorString("libvirt", "启动域失败", err.Error())
		return err
	}
	return nil
}

// ShutdownDomain 正常关闭指定的域
func ShutdownDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainShutdownFlags(domain, libvirt.DomainShutdownDefault)
	if err != nil {
		logger.ErrorString("libvirt", "关闭域失败", err.Error())
		return err
	}
	return nil
}

// ForceStopDomain 强制停止指定的域
func ForceStopDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainDestroyFlags(domain, libvirt.DomainDestroyDefault)
	if err != nil {
		logger.ErrorString("libvirt", "强制停止域失败", err.Error())
		return err
	}
	return nil
}

// RebootDomain 重启指定的域
func RebootDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainReboot(domain, libvirt.DomainRebootDefault)
	if err != nil {
		logger.ErrorString("libvirt", "重启域失败", err.Error())
		return err
	}
	return nil
}

// SuspendDomain 挂起指定的域
func SuspendDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainSuspend(domain)
	if err != nil {
		logger.ErrorString("libvirt", "挂起域失败", err.Error())
		return err
	}
	return nil
}

// ResumeDomain 恢复指定的域
func ResumeDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainResume(domain)
	if err != nil {
		logger.ErrorString("libvirt", "恢复域失败", err.Error())
		return err
	}
	return nil
}

// DeleteDomain 删除指定的域
func DeleteDomain(domain libvirt.Domain) error {
	err := Libvirt.DomainUndefine(domain)
	if err != nil {
		logger.ErrorString("libvirt", "删除域失败", err.Error())
		return err
	}
	return nil
}
