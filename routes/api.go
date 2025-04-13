// Package routes 注册路由
package routes

import (
	ctrl "gohub/app/http/controllers/api/v1"
	libCtrl "gohub/app/http/controllers/api/v1/libvirt"
	"gohub/app/http/middlewares"

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
		libGroup := v1.Group("/libvirt", middlewares.AuthApiKey())
		{
			libc := new(libCtrl.LibvirtController)
			{
				libGroup.GET("/domains", libc.ListAllDomains)
				libGroup.GET("/version", libc.GetLibVersion)
			}
		}
	}
}
