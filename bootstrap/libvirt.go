package bootstrap

import (
	"fmt"
	"govirt/pkg/config"
	"govirt/pkg/libvirtd"
	"govirt/pkg/logger"
	"govirt/pkg/network"
	"govirt/pkg/storagePool"
	"govirt/pkg/xmlDefine"
	"sync"
)

func InitLibvirt() {
	err := libvirtd.InitConnection(config.GetString("libvirt.conURI"))
	if err != nil {
		logger.FatalString("libvirt", "初始化libvirt连接", err.Error())
	}

	// 使用WaitGroup等待所有初始化任务完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 用于收集错误的channel
	errChan := make(chan error, 2)

	// 使用goroutine初始化存储池
	go func() {
		defer wg.Done()
		if err := InitStoragePool(); err != nil {
			errChan <- fmt.Errorf("初始化存储池失败: %w", err)
		}
	}()

	// 使用goroutine初始化网络
	go func() {
		defer wg.Done()
		if err := InitNetwork(); err != nil {
			errChan <- fmt.Errorf("初始化网络失败: %w", err)
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		logger.FatalString("libvirt", err.Error(), "")
	}

	logger.InfoString("libvirt", "初始化Libvirt控制器", "成功")
}

// InitStoragePool // InitStoragePool 初始化存储池
func InitStoragePool() error {
	var poolParams []xmlDefine.PoolTemplateParams
	// 定义需要初始化的存储池
	poolConfigs := []map[string]string{
		{"name": "pool.image.name", "path": "pool.image.path"},
		{"name": "pool.volume.name", "path": "pool.volume.path"},
	}

	// 构建参数列表
	for _, cfg := range poolConfigs {
		param := xmlDefine.PoolTemplateParams{
			Name: config.GetString(cfg["name"]),
			Path: config.GetString(cfg["path"]),
		}
		poolParams = append(poolParams, param)
	}

	// 批量初始化存储池
	return storagePool.InitSystemStoragePool(poolParams...)
}

func InitNetwork() error {
	var networkParams []xmlDefine.NetworkTemplateParams
	// 定义需要初始化的网络
	networkConfigs := []map[string]string{
		{"name": "network.internal.name",
			"domainName":  "network.internal.domainName",
			"forwardMode": "network.internal.forwardMode",
			"ip":          "network.internal.ip",
			"netmask":     "network.internal.netmask",
			"dhcpStart":   "network.internal.dhcpStart",
			"dhcpEnd":     "network.internal.dhcpEnd"},

		{"name": "network.external.name",
			"domainName":  "network.external.domainName",
			"forwardMode": "network.external.forwardMode",
			"ip":          "network.external.ip",
			"netmask":     "network.external.netmask",
			"dhcpStart":   "network.external.dhcpStart",
			"dhcpEnd":     "network.external.dhcpEnd"},
	}

	for _, cfg := range networkConfigs {
		param := xmlDefine.NetworkTemplateParams{
			Name:        config.GetString(cfg["name"]),
			DomainName:  config.GetString(cfg["domainName"]),
			ForwardMode: config.GetString(cfg["forwardMode"]),
			IPAddress:   config.GetString(cfg["ip"]),
			NetMask:     config.GetString(cfg["netmask"]),
			DhcpStart:   config.GetString(cfg["dhcpStart"]),
			DhcpEnd:     config.GetString(cfg["dhcpEnd"]),
		}
		networkParams = append(networkParams, param)
	}

	return network.InitSystemNetwork(networkParams...)
}
