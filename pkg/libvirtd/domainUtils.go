package libvirtd

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/xmlDefine"

	"github.com/digitalocean/go-libvirt"
)

// 定义DomainOperation类型
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

// ensureMacAddresses 确保 DomainTemplateParams 中的 MAC 地址已设置，如果未设置则生成
func (vc *VirtConn) ensureMacAddresses(dparams *xmlDefine.DomainTemplateParams) error {
	if dparams.InterMac == "" {
		macAddr, err := helpers.GenerateRandomMAC()
		if err != nil {
			return fmt.Errorf("生成内部网络MAC地址失败: %w", err)
		}
		dparams.InterMac = macAddr
	}
	if dparams.ExterMac == "" {
		macAddr, err := helpers.GenerateRandomMAC()
		if err != nil {
			return fmt.Errorf("生成外部网络MAC地址失败: %w", err)
		}
		dparams.ExterMac = macAddr
	}
	return nil
}
