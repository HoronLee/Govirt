// Package config 站点配置信息
package config

import "govirt/pkg/config"

// init 函数在包初始化时自动执行，用于设置应用的配置信息
func init() {
	// 调用 config 包的 Add 方法，添加一个名为 "app" 的配置项
	config.Add("app", func() map[string]any {
		// 返回一个包含应用配置信息的 map
		return map[string]any{

			// 应用名称，从环境变量 APP_NAME 中获取，默认值为 "govirt"
			"name": config.Env("APP_NAME", "govirt"),

			// 当前环境，用以区分多环境，一般为 local, stage, production, test
			"env": config.Env("APP_ENV", "production"),

			// 是否进入调试模式
			"debug": config.Env("APP_DEBUG", false),

			// 应用服务端口
			"port": config.Env("APP_PORT", "3000"),

			// 加密会话、JWT 加密
			"key": config.Env("APP_KEY", "114514horon1919810lee114514"),

			// 用以生成链接
			"url": config.Env("APP_URL", "http://localhost:3000"),

			// 设置时区，JWT 里会使用，日志记录里也会使用到
			"timezone": config.Env("TIMEZONE", "Asia/Shanghai"),
		}
	})
}
