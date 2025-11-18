package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// 测试上传图片接口
func TestAgentLogoUpload(t *testing.T) {
	target := "127.0.0.1:52320" //本机
	filePath := "C:\\Users\\Administrator\\Downloads\\133f23e2ce163748d3df547dbe3ecc5f.jpeg"
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(fmt.Errorf("无法打开文件: %v", err))
		return
	}
	defer file.Close()

	// 2. 创建multipart/form-data请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 3. 创建表单文件字段（字段名需与服务端一致，如"image"）
	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		t.Fatal(fmt.Errorf("创建表单文件字段失败: %v", err))
		return
	}

	// 4. 将文件内容复制到表单字段
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(fmt.Errorf("复制文件内容失败: %v", err))
		return
	}

	// 5. 关闭writer以完成multipart请求体的构建
	err = writer.Close()
	if err != nil {
		t.Fatal(fmt.Errorf("关闭multipart writer失败: %v", err))
		return
	}

	// 6. 创建HTTP POST请求
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/v1/agent/logo/upload", target), body)
	if err != nil {
		t.Fatal(fmt.Errorf("创建HTTP请求失败: %v", err))
		return
	}

	// 7. 设置请求头（Content-Type需包含boundary）
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 8. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(fmt.Errorf("发送请求失败: %v", err))
		return
	}

	// 9. 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(fmt.Errorf("读取响应失败: %v", err))
		return
	}

	fmt.Println(string(responseBody))
}
