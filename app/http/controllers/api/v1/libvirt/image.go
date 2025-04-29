package libvirt

import (
	"govirt/pkg/config"
	"govirt/pkg/libvirtd"
	"govirt/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ctrl *LibvirtController) CreateImageFromLocalFile(c *gin.Context) {
	var req struct {
		Name           string
		SourceFilePath string
		PoolName       string
		OSType         string
		Arch           string
		ImageType      string
		Description    string
		MinDisk        uint64
		MinRam         uint64
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, "请求参数错误")
		return
	}

	image, err := libvirtd.Conn.CreateImageFromLocalFile(
		req.Name,
		req.SourceFilePath,
		req.PoolName,
		req.OSType,
		req.Arch,
		req.ImageType,
		req.Description,
		req.MinDisk,
		req.MinRam,
	)
	if err != nil {
		response.Error(c, err, "创建镜像失败")
		return
	}

	response.Data(c, image)
}

func (ctrl *LibvirtController) DeleteImage(c *gin.Context) {
	ii := c.Query("image_identifier")

	if ii == "" {
		response.BadRequest(c, nil, "缺少镜像ID或UUID")
		return
	}

	err := libvirtd.Conn.DeleteImage(ii)
	if err != nil {
		response.Error(c, err, "删除镜像失败")
		return
	}

	response.Success(c)
}

func (ctrl *LibvirtController) ListActiveImages(c *gin.Context) {
	flag := c.DefaultQuery("flag", "0")
	flagInt, err := strconv.Atoi(flag)
	if err != nil {
		response.BadRequest(c, err, "flag参数无效")
		return
	}
	images, err := libvirtd.Conn.ListActiveImages(flagInt)
	if err != nil {
		response.Error(c, err, "获取活动镜像失败")
		return
	}

	response.Data(c, images)
}

func (ctrl *LibvirtController) SyncImages(c *gin.Context) {
	pi := c.DefaultQuery("pool_identifier", config.Get("pool.image.name"))
	err := libvirtd.Conn.SyncImagesWithVolumes(pi)
	if err != nil {
		response.Error(c, err, "同步镜像失败")
		return
	}
	response.Success(c)
}
