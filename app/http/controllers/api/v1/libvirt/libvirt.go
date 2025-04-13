package libvirt

import (
	"fmt"
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

func (ctrl *LibvirtController) ListAllDomains(c *gin.Context) {
	domains, err := libvirt.ListAllDomains()
	if err != nil {
		response.Error(c, err, "获取虚拟机列表失败")
		return
	}
	// 格式化UUID
	var formattedDomains []map[string]any
	for _, d := range domains {
		formattedDomains = append(formattedDomains, map[string]any{
			"Name": d.Name,
			"UUID": fmt.Sprintf("%x", d.UUID), // 转换为十六进制字符串
			"ID":   d.ID,
		})
	}

	response.Data(c, formattedDomains)
}
