package common

import (
	ccse2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ccse"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/crs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebm"
	ctebs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctebsbackup"
	ctecs2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctecs"
	sdkCtelb "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctelb"
	ctvpc2 "github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-core"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/amqp"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctebs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctecs"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctiam"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctimage"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/ctvpc"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mongodb"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/mysql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctyun-sdk-endpoint/pgsql"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/ctzos"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/dcs2"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/hpfs"
	ctgkafka "github.com/ctyun-it/terraform-provider-ctyun/internal/core/kafka"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/scaling"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/sfs"
	"sync"
)

var once sync.Once
var ctyunMetadata *CtyunMetadata

type CtyunMetadata struct {
	Apis          *Apis
	Credential    ctyunsdk.Credential
	extra         map[string]string
	SdkCredential core.Credential
}

// InitCtyunMetadata 初始化
func InitCtyunMetadata(apis *Apis, credential ctyunsdk.Credential, sdkCred core.Credential, extra map[string]string) {
	ctyunMetadata = &CtyunMetadata{Apis: apis, Credential: credential, SdkCredential: sdkCred, extra: extra}
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
	CtEbsApis       *ctebs.Apis
	CtEbsBackupApis *ctebsbackup.Apis
	CtEcsApis       *ctecs.Apis
	CtIamApis       *ctiam.Apis
	CtImageApis     *ctimage.Apis
	CtVpcApis       *ctvpc.Apis
	CtEbmApis       *ctebm.Apis
	SdkCtEbsApis    *ctebs2.Apis
	SdkCtEcsApis    *ctecs2.Apis
	SdkCtVpcApis    *ctvpc2.Apis
	SdkCtZosApis    *ctzos.Apis
	SdkCcseApis     *ccse2.Apis
	SdkDcs2Apis     *dcs2.Apis
	SdkCtElbApis    *sdkCtelb.Apis
	SdkCtMysqlApis  *mysql.Apis
	SdkCtPgsqlApis  *pgsql.Apis
	SdkKafkaApis    *ctgkafka.Apis
	SdkMongodbApis  *mongodb.Apis
	SdkAmqpApis     *amqp.Apis
	SdkCrsApis      *crs.Apis
	SdkHpfsApis     *hpfs.Apis
	SdkScalingApis  *scaling.Apis
	SdkSfsApi       *sfs.Apis
}
