package libvirt

import (
	v1 "govirt/app/http/controllers/api/v1"
	"govirt/pkg/libvirt"
	"govirt/pkg/response"

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
