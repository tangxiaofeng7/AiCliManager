package app

import "fmt"

// errInvalidParam 创建一个参数验证错误，用于 App 绑定层的输入校验
func errInvalidParam(msg string) error {
	return fmt.Errorf("参数错误: %s", msg)
}
