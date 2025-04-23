package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/response"

	"govirt/pkg/domain"

	"github.com/gin-gonic/gin"
)

// ListAllDomains 列出所有域
func (ctrl *LibvirtController) ListAllDomains(c *gin.Context) {
	domains, err := domain.ListAllDomains()
	if err != nil {
		response.Error(c, err, "列出所有域失败")
		return
	}
	response.Data(c, helpers.FormatStructSlice(domains))
}

// GetDomainState 获取指定域的状态
func (ctrl *LibvirtController) GetDomainStateByUUID(c *gin.Context) {
	uuid, err := helpers.UUIDStringToBytes(c.Query("uuid"))
	if err != nil {
		response.Error(c, err, "解析UUID失败")
		return
	}
	state, err := domain.GetDomainStateByUUID(uuid)
	if err != nil {
		response.Error(c, err, "获取域失败")
		return
	}
	response.Data(c, domain.DomainStateToString(state))
}

// UpdateDomainState 更新指定域的状态
func (ctrl *LibvirtController) UpdateDomainStateByUUID(c *gin.Context) {
	uuid, err := helpers.UUIDStringToBytes(c.Query("uuid"))
	if err != nil {
		response.Error(c, err)
		return
	}
	op := domain.StringToDomainOperation(c.Query("operation"))
	presentState, err := domain.UpdateDomainStateByUUID(uuid, op, 0)
	if err != nil {
		response.Error(c, err, "更新域状态失败")
		return
	}
	response.Data(c, domain.DomainStateToString(presentState))
}
