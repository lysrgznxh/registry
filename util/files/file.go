package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetFileName(fileName string) string {
	file := filepath.Base(fileName)
	strArr := strings.Split(file, ".")
	return strings.Join(strArr[:len(strArr)-1], ".")
}

// 读取文件内容
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err == nil {
		f.Close() //如果没有错误,关闭文件
	}
	return data, err
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func Trans2AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

// 去除多余空格
func StringCombineSpace(s string) string {
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(s, " ")
}

// 构建路径
func MakePath(paths ...string) string {
	path := ""
	for _, v := range paths {
		if path == "" {
			path = v
		} else {
			path = filepath.Join(path, v)
		}
	}
	return path
}

// 创建文件
func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

// 更新文件
func UpdateFile(path string, content string, flag int) error {
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	data := []byte(content)
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
