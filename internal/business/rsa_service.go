package business

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

// GetPublicKey 从Base64编码的公钥字符串解析出RSA公钥
func GetPublicKey(key string) (*rsa.PublicKey, error) {
	// 1. Base64解码
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %w", err)
	}

	// 2. 尝试直接解析为DER格式（兼容Java行为）
	if pubKey, err := x509.ParsePKIXPublicKey(keyBytes); err == nil {
		if rsaPub, ok := pubKey.(*rsa.PublicKey); ok {
			return rsaPub, nil
		}
	}

	// 3. 当DER解析失败时尝试PEM格式
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("无效的公钥格式：既不是DER也不是PEM")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %w", err)
	}

	if rsaPub, ok := pubKey.(*rsa.PublicKey); ok {
		return rsaPub, nil
	}

	return nil, errors.New("无法识别或不是RSA公钥")
}

func Encode(pw string) string {
	publicKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCCoIfL5oXcaYAmfo5Fpl3Fypo29JIhwuQTcmiwn/FBBW7L4L2Hm8rKbBpLZl/osCmoM6WSK8S9s/anp3ka1KIExUQDXYR8IH1cMse1+CzXcDZW/8Tp+NTPC/m6ymF00dUAi7k13q1Gv1Ez0jLx6ZqEu5tFpyVi0LRdLPc5G/T0DwIDAQAB"
	// 1. 解析PEM格式的公钥
	rsaPublicKey, err := GetPublicKey(publicKey)
	if err != nil {
		return ""
	}
	// 2. 要加密的原始数据
	plainText := []byte(pw)
	// 3. 使用公钥加密（选择PKCS#1 v1.5填充方案）
	// 对于更高安全性，建议使用rsa.EncryptOAEP
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, plainText)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}

	// 4. 输出Base64编码的密文（便于传输）
	encoded := base64.StdEncoding.EncodeToString(cipherText)
	//fmt.Println("加密结果 (Base64):")
	//fmt.Println(encoded)
	return encoded
}
