package goRuntime

import "os"

func IsExtraFile(filePath string) bool {
	// 文件不存在则返回error
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
