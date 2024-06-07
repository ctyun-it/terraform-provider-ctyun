package ctyunsdk

import (
	"errors"
	"os"
)

type Credential struct {
	ak string
	sk string
}

const (
	AkLen = 32
	SkLen = 32
)

// CheckAk 校验ak
func CheckAk(ak string) error {
	if len(ak) != AkLen {
		return errors.New("ak长度必须为32")
	}
	return nil
}

// CheckSk 校验sk
func CheckSk(sk string) error {
	if len(sk) != SkLen {
		return errors.New("sk长度必须为32")
	}
	return nil
}

// NewCredential 构造新的凭证
func NewCredential(ak string, sk string) (*Credential, error) {
	err := CheckAk(ak)
	if err != nil {
		return nil, err
	}

	err = CheckSk(sk)
	if err != nil {
		return nil, err
	}
	return &Credential{ak: ak, sk: sk}, nil
}

// NewCredentialFromEnv 构造新的凭证
func NewCredentialFromEnv() (*Credential, error) {
	ak := os.Getenv("CTYUN_AK")
	sk := os.Getenv("CTYUN_SK")
	return NewCredential(ak, sk)
}
