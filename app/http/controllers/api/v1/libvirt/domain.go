package libvirt

import (
	"fmt"
	"govirt/pkg/helpers"
	"govirt/pkg/logger"
	"govirt/pkg/response"
	"strconv"

	"govirt/pkg/domain"

	"github.com/gin-gonic/gin"
)

// ListAllDomains 列出所有域
func (ctrl *LibvirtController) ListAllDomains(c *gin.Context) {
	needResultInt, err := strconv.Atoi(c.DefaultQuery("needResults", "0"))
	if err != nil {
		logger.WarnString("libvirt", "ListAllDomains", fmt.Sprintf("resultNum转换失败: %s", err.Error()))
		needResultInt = -1
	}

	domains, err := domain.ListAllDomains(int32(needResultInt), 0)
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
