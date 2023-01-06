package util

import "os"

//判断文件是否已经生成. 如果controller和main已经生成了, 那就不再生成代码. 避免覆盖
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
