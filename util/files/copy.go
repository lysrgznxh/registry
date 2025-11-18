package files

import (
	"agent-network-protocol/registry/core"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 判断文件或目录是否存在
func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

// 拷贝文件
func CopyFile(src, dst string) (bool, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmNodeUtil)
	//dst := filepath.Join(dstDir, filepath.Base(src))
	log.Debug("dstPath:", dst)
	if len(src) == 0 || len(dst) == 0 {
		return false, nil
	}
	srcFile, e := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if e != nil {
		return false, e
	}
	defer srcFile.Close()

	dst = strings.Replace(dst, "\\", "/", -1)
	dstPathArr := strings.Split(dst, "/")
	dstPathArr = dstPathArr[0 : len(dstPathArr)-1]
	dstPath := strings.Join(dstPathArr, "/")

	dstFileInfo := GetFileInfo(dstPath)
	if dstFileInfo == nil {
		if e := os.MkdirAll(dstPath, os.ModePerm); e != nil {
			return false, e
		}
	}
	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	dstFile, e := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		return false, e
	}
	defer dstFile.Close()
	srcText, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return false, err
	}
	srcReader := bytes.NewBuffer(srcText) //如果文件属于脱敏白名单，则不进行脱敏
	if _, e := io.Copy(dstFile, srcReader); e != nil {
		return false, e
	} else {
		return true, nil
	}
}

// 拷贝目录
func CopyPath(src, dst string) (bool, error) {
	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil {
		return false, nil
	}
	if !srcFileInfo.IsDir() {
		return CopyFile(src, dst)
	}

	//创建目录
	dst = filepath.Join(dst, filepath.Base(src))
	_ = os.Mkdir(dst, os.ModePerm)

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relationPath := strings.Replace(path, src, "/", -1)
		dstPath := strings.TrimRight(strings.TrimRight(dst, "/"), "\\") + relationPath
		if !info.IsDir() {
			if _, err := CopyFile(path, dstPath); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			if _, err := os.Stat(dstPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
						return err
					} else {
						return nil
					}
				} else {
					return err
				}
			} else {
				return nil
			}
		}
	})

	if err != nil {
		return false, err
	}
	return true, nil
}
