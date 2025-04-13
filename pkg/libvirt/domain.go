package libvirt

import (
	"gohub/pkg/logger"

	"github.com/digitalocean/go-libvirt"
)

// ListAllDomains 连接到 libvirt 并列出所有域的信息
func ListAllDomains() ([]libvirt.Domain, error) {
	flags := libvirt.ConnectListDomainsActive | libvirt.ConnectListDomainsInactive
	domains, _, err := Libvirt.ConnectListAllDomains(1, flags)
	if err != nil {
		logger.ErrorString("libvirt", "列出所有域失败", err.Error())
		return nil, err
	}
	return domains, nil
}
