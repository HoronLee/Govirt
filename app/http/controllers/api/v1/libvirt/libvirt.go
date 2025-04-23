package libvirt

import (
	v1 "govirt/app/http/controllers/api/v1"
	"govirt/pkg/libvirtd"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

type LibvirtController struct {
	v1.BaseAPIController
}

func (ctrl *LibvirtController) GetServerInfo(c *gin.Context) {
	info, err := libvirtd.GetServerInfo()
	if err != nil {
		response.Error(c, err, "获取宿主机信息失败")
		return
	}
	response.Data(c, info)
}
