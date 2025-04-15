package libvirt

import (
	"govirt/pkg/libvirt"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAllNetworks 列出所有网络
func (ctrl *LibvirtController) ListAllNetworks(c *gin.Context) {
	networks, err := libvirt.ListAllNetworks()
	if err != nil {
		response.Error(c, err, "列出所有网络失败")
		return
	}
	response.Data(c, libvirt.FormatNetworks(networks))
}
