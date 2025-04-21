package helpers

import (
	"crypto/rand"
	"fmt"
	"sync"
)

var (
	// macLock 确保在生成MAC地址时的线程安全
	macLock sync.Mutex

	// 已经分配的MAC地址映射表，用于避免重复分配
	allocatedMACs = make(map[string]bool)
)

// GenerateRandomMAC 创建一个符合libvirt规范的随机MAC地址
// libvirt使用52:54:00作为默认的OUI前缀
func GenerateRandomMAC() (string, error) {
	macLock.Lock()
	defer macLock.Unlock()

	// 默认使用libvirt的OUI前缀52:54:00
	return GenerateRandomMACWithPrefix("52:54:00")
}

// GenerateRandomMACWithPrefix 使用指定的OUI前缀创建随机MAC地址
func GenerateRandomMACWithPrefix(prefix string) (string, error) {
	// 确保不生成重复的MAC地址
	for range 100 {
		// 生成3个随机字节
		b := make([]byte, 3)
		if _, err := rand.Read(b); err != nil {
			return "", fmt.Errorf("生成随机字节失败: %v", err)
		}

		// 构建MAC地址
		mac := fmt.Sprintf("%s:%02x:%02x:%02x", prefix, b[0], b[1], b[2])

		// 确保该地址未被使用
		if !allocatedMACs[mac] {
			allocatedMACs[mac] = true
			return mac, nil
		}
	}

	return "", fmt.Errorf("无法生成唯一的MAC地址，已达到最大重试次数")
}

// ReleaseMAC 释放一个已分配的MAC地址，允许将来重新使用
func ReleaseMAC(mac string) {
	macLock.Lock()
	defer macLock.Unlock()

	delete(allocatedMACs, mac)
}

// IsValidMAC 检查MAC地址的格式是否有效
func IsValidMAC(mac string) bool {
	// 简单的格式验证，可以使用更复杂的正则表达式
	_, err := fmt.Sscanf(mac, "%02x:%02x:%02x:%02x:%02x:%02x",
		new(int), new(int), new(int), new(int), new(int), new(int))
	return err == nil
}
