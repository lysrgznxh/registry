package api

import (
	"fmt"
	"io"

	"net/http"

	"testing"
)

// 测试列表接口
func TestAgentList(t *testing.T) {
	target := "127.0.0.1:52320" //本机

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/v1/agent/list", target), nil)
	if err != nil {
		t.Fatal(fmt.Errorf("创建HTTP请求失败: %v", err))
		return
	}

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
