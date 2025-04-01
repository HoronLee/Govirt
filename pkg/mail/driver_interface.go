package mail

type Driver interface {
	// 检查验证码
	// 原本返回的是bool，现在是error
	Send(email Email, config map[string]string) error
}
