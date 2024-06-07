package common

import (
	"sync"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	"terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
)

var once sync.Once
var ctyunMetadata *CtyunMetadata

type CtyunMetadata struct {
	Apis       *Apis
	Credential ctyunsdk.Credential
	extra      map[string]string
}

// InitCtyunMetadata 初始化
func InitCtyunMetadata(apis *Apis, credential ctyunsdk.Credential, extra map[string]string) {
	once.Do(func() {
		ctyunMetadata = &CtyunMetadata{Apis: apis, Credential: credential, extra: extra}
	})
}

// AcquireCtyunMetadata 获取实例对象
func AcquireCtyunMetadata() *CtyunMetadata {
	if ctyunMetadata == nil {
		panic("ctyun metadata not init!")
	}
	return ctyunMetadata
}

// GetExtra 获取默认设置的值
func (c CtyunMetadata) GetExtra(extraKey string) string {
	return c.extra[extraKey]
}

// GetExtraIfEmpty 如果目标值为空，获取默认设置的值，若目标值非空则返回目标值
func (c CtyunMetadata) GetExtraIfEmpty(target, extraKey string) string {
	if target == "" {
		return c.extra[extraKey]
	}
	return target
}

type Apis struct {
	CtEbsApis   *ctebs.Apis
	CtEcsApis   *ctecs.Apis
	CtIamApis   *ctiam.Apis
	CtImageApis *ctimage.Apis
	CtVpcApis   *ctvpc.Apis
}
