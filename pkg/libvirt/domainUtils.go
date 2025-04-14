package libvirt

import (
	"gohub/pkg/helpers"
	"gohub/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

type DomainOperation int32

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

// 操作类型到字符串的映射，将中文描述替换为英文描述
var domainOpStrings = map[DomainOperation]string{
	DomainOpStart:       "启动",
	DomainOpShutdown:    "关机",
	DomainOpForceStop:   "强制停止",
	DomainOpReboot:      "重启",
	DomainOpForceReboot: "强制重启",
	DomainOpSuspend:     "暂停",
	DomainOpResume:      "恢复",
	DomainOpSave:        "保存",
	DomainOpRestore:     "恢复状态",
	DomainOpDelete:      "删除",
	DomainOpClone:       "克隆",
	DomainOpMigrate:     "迁移",
	DomainOpSnapshot:    "快照",
	DomainOpUnknown:     "未知操作",
}

func (op DomainOperation) String() string {
	if s, ok := domainOpStrings[op]; ok {
		return s
	}
	return "未知操作"
}

// StringToDomainOperation 将字符串转换为DomainOperation
func StringToDomainOperation(s string) DomainOperation {
	switch s {
	case "Start":
		return DomainOpStart
	case "Shutdown":
		return DomainOpShutdown
	case "ForceStop":
		return DomainOpForceStop
	case "Reboot":
		return DomainOpReboot
	case "ForceReboot":
		return DomainOpForceReboot
	case "Suspend":
		return DomainOpSuspend
	case "Resume":
		return DomainOpResume
	case "Save":
		return DomainOpSave
	case "Restore":
		return DomainOpRestore
	case "Delete":
		return DomainOpDelete
	case "Clone":
		return DomainOpClone
	case "Migrate":
		return DomainOpMigrate
	case "Snapshot":
		return DomainOpSnapshot
	default:
		return DomainOpUnknown
	}
}

// FormatDomains 格式化域信息
func FormatDomains(domains []libvirt.Domain) []map[string]any {
	var formattedDomains []map[string]any
	for _, d := range domains {
		formattedDomains = append(formattedDomains, map[string]any{
			"ID":   d.ID,
			"Name": d.Name,
			"UUID": helpers.UUIDBytesToString(d.UUID),
		})
	}
	return formattedDomains
}

// GetDomainByUUID 根据 UUID 获取域
func GetDomainByUUID(uuid libvirt.UUID) (libvirt.Domain, error) {
	domains, err := ListAllDomains()
	if err != nil {
		logger.ErrorString("libvirt", "获取域失败", err.Error())
		return libvirt.Domain{}, err
	}
	for _, domain := range domains {
		if domain.UUID == uuid {
			return domain, nil
		}
	}
	return libvirt.Domain{}, nil
}

// DomainStateToString 将DomainState数值转换为对应的字符串状态
func DomainStateToString(state libvirt.DomainState) string {
	switch state {
	case libvirt.DomainNostate:
		return "Nostate"
	case libvirt.DomainRunning:
		return "Running"
	case libvirt.DomainBlocked:
		return "Blocked"
	case libvirt.DomainPaused:
		return "Paused"
	case libvirt.DomainShutdown:
		return "Shutdown"
	case libvirt.DomainShutoff:
		return "Shutoff"
	case libvirt.DomainCrashed:
		return "Crashed"
	case libvirt.DomainPmsuspended:
		return "Pmsuspended"
	default:
		return "Unknown"
	}
}

// DomainJobOperationToString 将DomainJobOperation数值转换为对应的字符串
func DomainJobOperationToString(op libvirt.DomainJobOperation) string {
	switch op {
	case libvirt.DomainJobOperationStrStart:
		return "Start"
	case libvirt.DomainJobOperationStrSave:
		return "Save"
	case libvirt.DomainJobOperationStrRestore:
		return "Restore"
	case libvirt.DomainJobOperationStrMigrationIn:
		return "MigrationIn"
	case libvirt.DomainJobOperationStrMigrationOut:
		return "MigrationOut"
	case libvirt.DomainJobOperationStrSnapshot:
		return "Snapshot"
	case libvirt.DomainJobOperationStrSnapshotRevert:
		return "SnapshotRevert"
	case libvirt.DomainJobOperationStrDump:
		return "Dump"
	case libvirt.DomainJobOperationStrBackup:
		return "Backup"
	default:
		return "Unknown"
	}
}
