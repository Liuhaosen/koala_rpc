package errno

import "errors"

var (
	NotHaveInstance  = errors.New("没有找到实例")
	ErrNotHaveNodes  = errors.New("没有节点")
	ErrAllNodeFailed = errors.New("所有节点尝试失败")
	ErrConnectError  = errors.New("连接失败")
	InvalidNode      = errors.New("节点不存在")
	ConnFailed       = errors.New("连接失败")
)

func ISConnectError(err error) bool {
	return err == ErrConnectError
}
