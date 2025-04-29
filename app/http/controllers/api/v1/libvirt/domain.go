package libvirt

import (
	"govirt/pkg/helpers"
	"govirt/pkg/libvirtd"
	"govirt/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAllDomains 列出所有域
func (ctrl *LibvirtController) ListAllDomains(c *gin.Context) {
	domains, err := libvirtd.Conn.ListAllDomains()
	if err != nil {
		response.Error(c, err, "列出所有域失败")
		return
	}
	response.Data(c, helpers.FormatStructSlice(domains))
}

// GetDomainState 获取指定域的状态
func (ctrl *LibvirtController) GetDomainStateByUUID(c *gin.Context) {
	di, err := helpers.UUIDStringToBytes(c.Query("domain_identifier"))
	if err != nil {
		response.Error(c, err, "解析UUID失败")
		return
	}

	domain, err := libvirtd.Conn.GetDomain(di)
	if err != nil {
		response.Error(c, err, "获取域失败")
		return
	}

	state, err := libvirtd.Conn.GetDomainState(domain)
	if err != nil {
		response.Error(c, err, "获取域状态失败")
		return
	}

	response.Data(c, libvirtd.DomainStateToString(state))
}

// UpdateDomainState 更新指定域的状态
func (ctrl *LibvirtController) UpdateDomainStateByUUID(c *gin.Context) {
	di, err := helpers.UUIDStringToBytes(c.Query("domain_identifier"))
	if err != nil {
		response.Error(c, err)
		return
	}

	op := libvirtd.StringToDomainOperation(c.Query("operation"))
	presentState, err := libvirtd.Conn.UpdateDomainStateByUUID(di, op, 0)
	if err != nil {
		response.Error(c, err, "更新域状态失败")
		return
	}

	response.Data(c, libvirtd.DomainStateToString(presentState))
}

// DeleteDomain 删除指定域
func (ctrl *LibvirtController) DeleteDomain(c *gin.Context) {
	di, err := helpers.UUIDStringToBytes(c.Query("domain_identifier"))
	if err != nil {
		response.Error(c, err)
		return
	}

	dm, err := libvirtd.Conn.GetDomain(di)
	if err != nil {
		response.Error(c, err, "获取域失败")
		return
	}

	err = libvirtd.Conn.DeleteStoppedDomain(dm)
	if err != nil {
		response.Error(c, err, "删除域失败")
		return
	}

	response.Success(c)
}
