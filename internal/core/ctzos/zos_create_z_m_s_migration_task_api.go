package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosCreateZMSMigrationTaskApi
/* 创建对象存储迁移任务
 */type ZosCreateZMSMigrationTaskApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosCreateZMSMigrationTaskApi(client *core.CtyunClient) *ZosCreateZMSMigrationTaskApi {
	return &ZosCreateZMSMigrationTaskApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodPost,
			UrlPath:      "/v4/zms/create-migration",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosCreateZMSMigrationTaskApi) Do(ctx context.Context, credential core.Credential, req *ZosCreateZMSMigrationTaskRequest) (*ZosCreateZMSMigrationTaskResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	_, err := ctReq.WriteJson(req, a.template.ContentType)
	if err != nil {
		return nil, err
	}
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosCreateZMSMigrationTaskResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosCreateZMSMigrationTaskRequest struct {
	RegionID         string `json:"regionID,omitempty"`         /*  资源池 ID  */
	MigrationName    string `json:"migrationName,omitempty"`    /*  任务名称，必须为大小写字母、数字、横线或下划线，长度在4-32个字符之间，且名称不能重复  */
	StorageType      string `json:"storageType,omitempty"`      /*  迁移到目的端存储类型，默认为标准存储。①MATCH_RESOURCE：匹配源端，匹配源端存储类型时，仅能自动匹配源端的“标准”和“低频”类型；匹配源端的“归档”或“深度归档”类型，请您务必提前对源端归档数据进行手动解冻，并确保迁移任务完成前数据保持解冻状态，否则该部分数据会迁移失败；②STANDARD：标准存储；③STANDARD_IA：低频存储；④GLACIER：归档存储  */
	AclConf          string `json:"aclConf,omitempty"`          /*  目的端ACL配置，默认为匹配源端，①match-resource：匹配源端；②private：私有；③public-read：公共读  */
	ConflictMode     string `json:"conflictMode,omitempty"`     /*  同名文件处理选项，默认为IGNORE,①OVERWRITE：同名文件进行覆盖；                                       ②IGNORE：同名文件进行忽略；③COMPARE：同名文件按最后修改时间(即LastModified)比较，若源端LastModified小于目的端LastModified，则此文件被执行跳过；若源端LastModified大于目的端LastModified，则执行覆盖；若源端与目的端文件LastModified一致，则判断两者的文件大小，大小一致则执行跳过，大小不一致则执行覆盖。  */
	MigrateStartTime string `json:"migrateStartTime,omitempty"` /*  迁移晚于起始时间的对象，该选项会迁移最后修改时间(即LastModified)晚于指定时间的对象。可以设置两种格式"year-month-day hour:minute:second"或"year-month-day"。可填时间范围限制为[1970-01-02 00:00:00,2037-12-31 23:59:59]，若同时填入migrateStartTime和migrateEndTime，则migrateStartTime值应小于migrateEndTime。  */
	MigrateEndTime   string `json:"migrateEndTime,omitempty"`   /*  迁移早于终止时间的对象，该选项会迁移最后修改时间(即LastModified)早于指定时间的对象。可以设置两种格式"year-month-day hour:minute:second"或"year-month-day"，默认为当前任务创建时间加10年。。可填时间范围限制为[1970-01-02 00:00:00,2037-12-31 23:59:59]，若同时填入migrateStartTime和migrateEndTime，则migrateStartTime值应小于migrateEndTime。  */
}

type ZosCreateZMSMigrationTaskResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
