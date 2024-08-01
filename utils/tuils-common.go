package uutils

import (
	"os"
)

/**
 * @Description: 删除文件夹中的所有内容
 * @param path 文件夹路径
 * @return error
 */
func RemoveAll(path string) error {
	// 删除文件夹中的所有内容
	/* 	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
	   		if err != nil {
	   			return err
	   		}
	   		if info.IsDir() {
	   			return nil // 继续递归进入子目录
	   		}
	   		return os.Remove(path) // 删除文件
	   	})
	   	if err != nil {
	   		return err
	   	} */
	// 删除文件夹
	return os.RemoveAll(path)
}
