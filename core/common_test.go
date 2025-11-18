package core

import (
	"agent-network-protocol/registry/util/crypto"
	"fmt"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"os"
	"testing"
)

func TestVersionCheck(t *testing.T) {
	testCases := []struct {
		version string
		valid   bool
	}{
		{"0.0.0", true},    // 符合格式
		{"0.0.1", false},   // 最后一段非 0
		{"0.1.0", false},   // 第二段非 0
		{"1.0.0", false},   // 第一段非 0
		{"0.0", false},     // 缺少第三段
		{"0.0.0.0", false}, // 多出一段
		{"0.0.a", false},   // 包含非数字字符
		{"0.0.0-", false},  // 包含非法字符
		{"0.0.0+", false},  // 包含非法字符
	}

	// 执行测试
	for _, tc := range testCases {
		result := IsValidVersion(tc.version)
		fmt.Printf("版本号: %-8s 预期结果: %-5v 实际结果: %-5v 通过: %v\n",
			tc.version, tc.valid, result, result == tc.valid)
	}
}

func TestNewPrivateKey(t *testing.T) {
	key, err := crypto.GenerateECKeyPair(crypto.Secp256k1())
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

// 测试注册签名
func TestMakeRegisterSign(t *testing.T) {
	content := "hello world"
	privateKeyPath := "D:\\test-private-key.pem"
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatal(fmt.Errorf("read private key: %w", err))
		return
	}
	key, err := crypto.PrivateKeyFromPEM(keyBytes)
	if err != nil {
		t.Fatal(fmt.Errorf("decode private key: %w", err))
		return
	}
	signByData := []byte(content)
	signBytes, err := crypto.MakeSignature(key, signByData)
	if err != nil {
		t.Fatal(fmt.Errorf("sign payload: %w", err))
		return
	}
	address := ethcrypto.PubkeyToAddress(key.PublicKey).Hex()
	t.Log("签名字符串:", content)
	t.Log("对应的以太坊地址: ", address)
	t.Log("签名:", signBytes)
	return
}

// 验证注册签名
func TestCheckRegisterSign(t *testing.T) {
	privateKeyPath := "D:\\test-private-key.pem"
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatal(fmt.Errorf("read private key: %w", err))
		return
	}
	privateKey, err := crypto.PrivateKeyFromPEM(keyBytes)
	if err != nil {
		t.Fatal(fmt.Errorf("decode private key: %w", err))
		return
	}
	content := []byte("hello world")
	signStr := "UQ0kzU2vK4W705SwtKMzI_zy1IVH1rKuqibknLzfS1PO8RLKzx5O5h0yEZDutiMSQRXtwr2dr72hx5IV8uvMmw"
	t.Log("sign:", signStr)

	pubkeyAddr := ethcrypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	t.Log("公钥地址:", pubkeyAddr)

	ok := crypto.VerifySignature(content, signStr, pubkeyAddr)

	if ok {
		t.Log("签名正确")
	} else {
		t.Log("签名错误")
	}
}
