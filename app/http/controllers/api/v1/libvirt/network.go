package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/network"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAllNetworks 列出所有网络
func (ctrl *LibvirtController) ListAllNetworks(c *gin.Context) {
	networks, err := network.ListAllNetworks()
	if err != nil {
		response.Error(c, err, "列出所有网络失败")
		return
	}
	response.Data(c, helpers.FormatStructSlice(networks))
}
