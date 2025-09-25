package ctzos

import (
	"context"
	"github.com/ctyun-it/terraform-provider-ctyun/internal/core/core"
	"net/http"
)

// ZosGetZMSMigrationTaskDetailApi
/* 查询对象存储迁移任务详情
 */type ZosGetZMSMigrationTaskDetailApi struct {
	template core.CtyunRequestTemplate
	client   *core.CtyunClient
}

func NewZosGetZMSMigrationTaskDetailApi(client *core.CtyunClient) *ZosGetZMSMigrationTaskDetailApi {
	return &ZosGetZMSMigrationTaskDetailApi{
		client: client,
		template: core.CtyunRequestTemplate{
			EndpointName: EndpointName,
			Method:       http.MethodGet,
			UrlPath:      "/v4/zms/get-migration-detail",
			ContentType:  "application/json",
		},
	}
}

func (a *ZosGetZMSMigrationTaskDetailApi) Do(ctx context.Context, credential core.Credential, req *ZosGetZMSMigrationTaskDetailRequest) (*ZosGetZMSMigrationTaskDetailResponse, error) {
	builder := core.NewCtyunRequestBuilder(a.template)
	builder.WithCredential(credential)
	ctReq := builder.Build()
	ctReq.AddParam("regionID", req.RegionID)
	ctReq.AddParam("migrationID", req.MigrationID)
	response, err := a.client.RequestToEndpoint(ctx, ctReq)
	if err != nil {
		return nil, err
	}
	var resp ZosGetZMSMigrationTaskDetailResponse
	err = response.Parse(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type ZosGetZMSMigrationTaskDetailRequest struct {
	RegionID    string /*  资源池 ID  */
	MigrationID string /*  迁移任务ID  */
}

type ZosGetZMSMigrationTaskDetailResponse struct {
	StatusCode  int32  `json:"statusCode,omitempty"`  /*  返回状态码(800为成功，900为处理中/失败)  */
	Message     string `json:"message,omitempty"`     /*  状态描述  */
	Description string `json:"description,omitempty"` /*  状态描述，一般为中文  */
	ErrorCode   string `json:"errorCode,omitempty"`   /*  业务细分码（仅失败时具有此参数），为product.module.code三段式码  */
	Error       string `json:"error,omitempty"`       /*  业务细分码（大驼峰形式，仅失败时具有此参数），为 Product.Module.Code 三段式码  */
}
