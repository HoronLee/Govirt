// Package routes 注册路由
package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		akGroup := v1.Group("/api", middlewares.AuthApiKey())
		{
			apic := new(controllers.ApikeyController)
			{
				akGroup.GET("", apic.ListApikey)
				akGroup.POST("", apic.CreateApikey)
				akGroup.DELETE("/:name", apic.DeleteApikey)
			}
		}
	}
}
