package core

import (
	"os"
)

type Credential struct {
	ak string
	sk string
}

// NewCredential 构造新的凭证
func NewCredential(ak string, sk string) *Credential {
	return &Credential{ak: ak, sk: sk}
}

// CredentialFromEnv 构造新的凭证
func CredentialFromEnv() *Credential {
	ak := os.Getenv("CTYUN_AK")
	sk := os.Getenv("CTYUN_SK")
	return NewCredential(ak, sk)
}

// GetAccessKey 获取 Access Key
func (c *Credential) GetAccessKey() string {
	return c.ak
}

// GetSecretKey 获取 Secret Key
func (c *Credential) GetSecretKey() string {
	return c.sk
}
