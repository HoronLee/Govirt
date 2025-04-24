package domain

import (
	"fmt"
	"govirt/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// DomainOperation 定义DomainOperation类型
type DomainOperation int32

// 定义DomainOperation常量
const (
	DomainOpStart       DomainOperation = iota + 1 // 启动
	DomainOpShutdown                               // 正常关机
	DomainOpForceStop                              // 强制停止
	DomainOpReboot                                 // 重启
	DomainOpForceReboot                            // 强制重启
	DomainOpSuspend                                // 暂停
	DomainOpResume                                 // 恢复
	DomainOpSave                                   // 保存状态
	DomainOpRestore                                // 恢复状态
	DomainOpDelete                                 // 删除
	DomainOpClone                                  // 克隆
	DomainOpMigrate                                // 迁移
	DomainOpSnapshot                               // 创建快照
	DomainOpUnknown                                // 未知操作
)

// stringToDomainOp 用于将字符串映射到DomainOperation
var stringToDomainOp = map[string]DomainOperation{
	"Start":       DomainOpStart,
	"Shutdown":    DomainOpShutdown,
	"ForceStop":   DomainOpForceStop,
	"Reboot":      DomainOpReboot,
	"ForceReboot": DomainOpForceReboot,
	"Suspend":     DomainOpSuspend,
	"Resume":      DomainOpResume,
	"Save":        DomainOpSave,
	"Restore":     DomainOpRestore,
	"Delete":      DomainOpDelete,
	"Clone":       DomainOpClone,
	"Migrate":     DomainOpMigrate,
	"Snapshot":    DomainOpSnapshot,
}

// StringToDomainOperation 将字符串转换为DomainOperation
func StringToDomainOperation(s string) DomainOperation {
	if op, ok := stringToDomainOp[s]; ok {
		return op
	}
	return DomainOpUnknown
}

var domainStateToString = map[libvirt.DomainState]string{
	libvirt.DomainNostate:     "Nostate",
	libvirt.DomainRunning:     "Running",
	libvirt.DomainBlocked:     "Blocked",
	libvirt.DomainPaused:      "Paused",
	libvirt.DomainShutdown:    "Shutdown",
	libvirt.DomainShutoff:     "Shutoff",
	libvirt.DomainCrashed:     "Crashed",
	libvirt.DomainPmsuspended: "Pmsuspended",
}

// DomainStateToString 将DomainState数值转换为对应的字符串状态
func DomainStateToString(state libvirt.DomainState) string {
	if str, ok := domainStateToString[state]; ok {
		return str
	}
	return "Unknown"
}

// GetDomain 根据 UUID 或名称获取域
func GetDomain(identifier any) (libvirt.Domain, error) {
	domains, err := ListAllDomains(-1, 1|2)
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
			if domain.Name == id {
				return domain, nil
			}
		default:
			return libvirt.Domain{}, fmt.Errorf("无效的标识符类型: %T", identifier)
		}
	}

	return libvirt.Domain{}, nil
}
