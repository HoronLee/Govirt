// Package routes 注册路由
package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册网页相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			// 注册
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", suc.IsPhoneExist)
			// 判断 Email 是否已注册
			authGroup.POST("/signup/email/exist", suc.IsEmailExist)
			// 使用 Email 注册
			authGroup.POST("/signup/using-email", suc.SignupUsingEmail)
			// ------------

			// 验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", vcc.ShowCaptcha)
			// 发送邮件验证码
			authGroup.POST("/verify-codes/email", vcc.SendUsingEmail)
			// ------------

			// 登录
			lgc := new(auth.LoginController)
			// 使用手机号，短信验证码进行登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// 使用密码登录，支持手机号，Email 和 用户名
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			// 刷新 Access Token
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)
			// ------------

			// 密码
			pwc := new(auth.PasswordController)
			// 使用 Email 重制密码
			authGroup.POST("/password-reset/using-email", pwc.ResetByEmail)
			// ------------

			// 用户
			uc := new(controllers.UsersController)
			// 获取当前用户
			v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
			// ------
		}
	}
}
