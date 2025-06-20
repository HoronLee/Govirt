// Package routes 注册路由
package routes

import (
	ctrl "govirt/app/http/controllers/api/v1"
	libCtrl "govirt/app/http/controllers/api/v1/libvirt"
	"govirt/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		akGroup := v1.Group("/api", middlewares.AuthApiKey())
		{
			apic := new(ctrl.ApikeyController)
			{
				akGroup.GET("", apic.ListApikey)
				akGroup.POST("", apic.CreateApikey)
				akGroup.DELETE("/:name", apic.DeleteApikey)
			}
		}
		// libvirt 相关路由
		libGroup := v1.Group("/libvirt", middlewares.AuthApiKey())
		{
			libc := new(libCtrl.LibvirtController)
			{
				libGroup.GET("/info", libc.GetServerInfo)
				// domain 相关路由
				domainGroup := libGroup.Group("/domain")
				{
					domainGroup.GET("/all", libc.ListAllDomains)
					domainGroup.GET("/state", libc.GetDomainStateByUUID)
					domainGroup.PUT("/state", libc.UpdateDomainStateByUUID)
					domainGroup.DELETE("/delete", libc.DeleteDomain)
				}
				// network 相关路由
				networkGroup := libGroup.Group("/network")
				{
					networkGroup.GET("/all", libc.ListAllNetworks)
					networkGroup.POST("/createStart", libc.CreateNetwork)
					networkGroup.DELETE("/stopDelete", libc.DeleteNetwork)
					networkGroup.POST("/active", libc.ActiveNetwork)
				}
				// storagePool 相关路由
				poolGroup := libGroup.Group("/pool")
				{
					poolGroup.GET("/all", libc.ListAllStoragePools)
					poolGroup.POST("/createStart", libc.CreateStartStoragePool)
					poolGroup.DELETE("/stopDelete", libc.DeleteStoragePool)
				}
				// volume 相关路由
				volumeGroup := libGroup.Group("/volume")
				{
					volumeGroup.GET("/allSummary", libc.ListVolumesSummary)
					volumeGroup.GET("/allDetail", libc.ListVolumesDetails)
					volumeGroup.POST("/create", libc.CreateVolume)
					volumeGroup.DELETE("/delete", libc.DeleteVolume)
					volumeGroup.POST("/clone", libc.CloneVolume)
				}
				// image 相关路由
				imageGroup := libGroup.Group("/image")
				{
					imageGroup.POST("/create", libc.CreateImageFromLocalFile)
					imageGroup.DELETE("/delete", libc.DeleteImage)
					imageGroup.GET("/listActive", libc.ListActiveImages)
					imageGroup.POST("/sync", libc.SyncImages)
				}
			}
		}
	}
}
