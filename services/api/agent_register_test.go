package api

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/core/model"
	"agent-network-protocol/registry/util/crypto"
	"agent-network-protocol/registry/util/net"
	"agent-network-protocol/registry/util/serializer"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestNewPrivateKey(t *testing.T) {
	key, err := newPrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}
	fileContent, err := crypto.PrivateKeyToPEM(key)
	if err != nil {
		t.Fatal(err)
		return
	}
	filepath := "D:\\test-private-key.pem"
	err = os.WriteFile(filepath, fileContent, 0600)
	if err != nil {
		t.Fatal("保存密钥文件失败:", err)
	}
	t.Log("密钥保证成功:", filepath)
}

// 创建新私钥
func newPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateECKeyPair(crypto.Secp256k1())
}

// 从文件加载私钥
func loadPrivateKey() (*ecdsa.PrivateKey, error) {
	privateKeyPath := "D:\\test-private-key.pem"
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := crypto.PrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// 测试调用接口,登记agent信息
func TestAgentRegister(t *testing.T) {
	var target = "127.0.0.1:52320" // 本机
	ctx := context.Background()
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: time.Second * 60, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	//key, err := newPrivateKey()
	key, err := loadPrivateKey()
	if err != nil {
		t.Fatal(err)
		return
	}
	name := "my.agent"
	title := "我的智能体2"
	nonce := uuid.New().String()
	signByData := []byte(fmt.Sprintf("%s_%s", name, nonce))
	signString, err := crypto.MakeSignature(key, signByData)
	if err != nil {
		t.Fatal(fmt.Errorf("sign payload: %w", err))
		return
	}
	address := ethcrypto.PubkeyToAddress(key.PublicKey).Hex()

	var req = core.AgentRegisterReq{
		Sign:        signString,
		Nonce:       nonce,
		Title:       title,
		PublicAddr:  address,
		Author:      "auth",
		Category:    core.CategroyYingXiao,
		Name:        name,
		Version:     "1.0.1",
		Type:        model.AgentTypeAgent,
		Logo:        "",
		Description: "智能体的能力描述",
		Status:      model.StatusDeleted,
		Remotes: []model.Transport{
			{
				Type: "streamable-http",
				URL:  "https://127.0.0.1/v1/chat-messages",
				Headers: []model.KeyValueInput{
					{
						InputWithVariables: model.InputWithVariables{
							Input: model.Input{
								Description: "api key",
								IsRequired:  true,
								Format:      model.FormatString,
								Value:       "Bearer ",
								Default:     "Bearer ",
							},
						},
						Name: "Authorization",
					},
				},
				Jsons: []model.Argument{
					{
						Name:       "query",
						IsRepeated: false,
						InputWithVariables: model.InputWithVariables{
							Input: model.Input{
								Description: "要查询的内容",
								IsRequired:  true,
								Format:      model.FormatString,
							},
						},
					},
					{
						Name:       "inputs",
						IsRepeated: false,
						InputWithVariables: model.InputWithVariables{
							Input: model.Input{
								Description: "用户自定义结构体参数",
								IsRequired:  false,
								IsArray:     false,
								Format:      model.FormatJson,
							},
							Variables: map[string]model.Input{
								"variable_name1": {
									IsRequired: false,
									IsArray:    true,
									Format:     model.FormatJson,
									Variables: map[string]model.Input{
										"type": {
											Description: "文件类型",
											IsRequired:  false,
											Format:      model.FormatString,
											Choices:     []string{"image", "document", "audio", "video", "custom"},
										},
										"transfer_method": {
											Description: "传输方法",
											IsRequired:  false,
											Format:      model.FormatString,
											Choices:     []string{"remote_url", "local_file"},
										},
										"url": {
											Description: "图片地址",
											Format:      model.FormatFilePath,
										},
										"upload_file_id": {
											Description: "上传文件id",
											Format:      model.FormatString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	res := &serializer.ApiReturnInterface{}

	err = net.PostJSON(ctx, &client, fmt.Sprintf("http://%s/v1/agent/reg", target), req, res)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(res)

	xxx, _ := json.Marshal(res)
	fmt.Println(string(xxx))

	if res.Status != 1 {
		fmt.Println(res.Info)
		return
	}
	t.Log(res)

}

func TestNameGrep(t *testing.T) {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{3,50}$`)
	testStrings := []string{
		"Hello_World123.txtHello_World123.txtHello_World123.txtorld123.txtHello_World123.txtorld123.txtHello_World123.txt",
		"Hello_World123.txtHello_World123.txtHello_World123.txt",
		"Hello_World123.txt",
		"valid..example",
		"invalid@email.com",
		"another-invalid",
		"a",
		"ab",
		"abc",
		"", // 空字符串，因为使用了*，所以也是允许的
	}

	for _, str := range testStrings {
		// 使用 MatchString 方法进行匹配
		if regex.MatchString(str) {
			fmt.Printf("'%s' -> 符合要求\n", str)
		} else {
			fmt.Printf("'%s' -> 不符合要求\n", str)
		}
	}
}
