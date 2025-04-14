package libvirt

import (
	v1 "gohub/app/http/controllers/api/v1"
	"gohub/pkg/libvirt"
	"gohub/pkg/response"

	"github.com/gin-gonic/gin"
)

type LibvirtController struct {
	v1.BaseAPIController
}

func (ctrl *LibvirtController) GetLibVersion(c *gin.Context) {
	version, err := libvirt.GetLibVersion()
	if err != nil {
		response.Error(c, err, "获取libvirt版本失败")
		return
	}
	response.Data(c, version)

}
