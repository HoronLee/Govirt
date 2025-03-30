// Package app 应用信息
// 这类的操作属于应用级别的，所以将其他们放置于自建的 app 包里
package app

import (
	"gohub/pkg/config"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}
